package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"

	"amap"

	"gopkg.in/yaml.v3"
)

//go:embed index.html
var indexHTML string

type config struct {
	JSAPI struct {
		Key         string `yaml:"key"`
		SecurityKey string `yaml:"security_key"`
	} `yaml:"js_api"`
	WebAPI struct {
		Key string `yaml:"key"`
	} `yaml:"web_api"`
	Destinations []string `yaml:"destinations"`
}

type destInfo struct {
	Address    string `json:"address"`
	Coordinate string `json:"coordinate"`
	Lng        string `json:"lng"`
	Lat        string `json:"lat"`
}

type pageData struct {
	JSKey       string
	SecurityKey string
	DestsJSON   string
}

type destResult struct {
	Destination string `json:"destination"`
	Duration    string `json:"duration"`
	Distance    string `json:"distance"`
	Seconds     int64  `json:"seconds"`
	Meters      int64  `json:"meters"`
	Error       string `json:"error,omitempty"`
}

var destinations []destInfo

func main() {
	raw, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("读取 config.yaml 失败: %v", err)
	}
	var cfg config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		log.Fatalf("解析 config.yaml 失败: %v", err)
	}
	if len(cfg.Destinations) == 0 {
		log.Fatal("config.yaml 中未配置 destinations")
	}

	amap.AmapKey = cfg.WebAPI.Key

	for _, addr := range cfg.Destinations {
		fmt.Printf("正在解析目的地: %s ... ", addr)
		geocodes, err := amap.GeoCode(addr, "")
		if err != nil {
			log.Fatalf("\n地理编码 [%s] 失败: %v", addr, err)
		}
		loc := geocodes[0].Location
		parts := strings.Split(loc, ",")
		if len(parts) != 2 {
			log.Fatalf("\n坐标格式异常: %s", loc)
		}
		fmt.Printf("%s\n", loc)
		destinations = append(destinations, destInfo{
			Address: addr, Coordinate: loc, Lng: parts[0], Lat: parts[1],
		})
	}

	destsJSON, _ := json.Marshal(destinations)
	pd := pageData{
		JSKey:       cfg.JSAPI.Key,
		SecurityKey: cfg.JSAPI.SecurityKey,
		DestsJSON:   string(destsJSON),
	}

	tmpl := template.Must(template.New("index").Parse(indexHTML))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, pd)
	})
	http.HandleFunc("/api/driving", handleDriving)
	http.HandleFunc("/api/batch-driving", handleBatchDriving)
	http.HandleFunc("/api/results", handleResults)

	addr := ":6052"
	fmt.Printf("服务启动: http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// ---------- 单点点击 ----------

func handleDriving(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	if origin == "" {
		writeJSON(w, map[string]string{"error": "缺少 origin 参数"})
		return
	}

	address := origin
	if regeo, err := amap.ReGeoCode(origin); err == nil {
		address = regeo.FormattedAddress
	}

	results := make([]destResult, len(destinations))
	var wg sync.WaitGroup
	for i, dest := range destinations {
		wg.Add(1)
		go func(idx int, d destInfo) {
			defer wg.Done()
			dr := destResult{Destination: d.Address}
			if detail, err := amap.GetDrivingDetail(origin, d.Coordinate); err == nil {
				dr.Duration = formatDuration(detail.Duration)
				dr.Distance = formatDistance(detail.DistanceMeters)
				dr.Seconds = int64(detail.Duration.Seconds())
				dr.Meters = detail.DistanceMeters
			} else {
				dr.Error = err.Error()
			}
			results[idx] = dr
		}(i, dest)
	}
	wg.Wait()

	var dests []fileDestRecord
	for i, dr := range results {
		if dr.Error != "" {
			continue
		}
		dests = append(dests, fileDestRecord{
			Address: destinations[i].Address, Coordinate: destinations[i].Coordinate,
			Duration: dr.Duration, Seconds: dr.Seconds,
			Distance: dr.Distance, Meters: dr.Meters,
		})
	}
	if len(dests) > 0 {
		appendResultsToFile([]filePointRecord{{
			Address: address, Coordinate: origin, Destinations: dests,
		}})
	}

	writeJSON(w, map[string]interface{}{"address": address, "results": results})
}

// ---------- 批量画圈 ----------

type batchRequest struct {
	Points     []string `json:"points"`
	IntervalMs int      `json:"interval_ms"`
}

type batchPointResult struct {
	Index        int          `json:"index"`
	Origin       string       `json:"origin"`
	Address      string       `json:"address"`
	Destinations []destResult `json:"destinations"`
}

func handleBatchDriving(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req batchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": "无效的请求体"})
		return
	}
	if len(req.Points) == 0 {
		writeJSON(w, map[string]string{"error": "没有提供坐标点"})
		return
	}

	interval := time.Duration(req.IntervalMs) * time.Millisecond
	results := make([]batchPointResult, len(req.Points))
	callCount := 0

	for i, point := range req.Points {
		item := batchPointResult{Index: i + 1, Origin: point}

		if callCount > 0 && interval > 0 {
			time.Sleep(interval)
		}
		callCount++
		if regeo, err := amap.ReGeoCode(point); err == nil {
			item.Address = regeo.FormattedAddress
		} else {
			item.Address = point
		}

		item.Destinations = make([]destResult, len(destinations))
		for j, dest := range destinations {
			if interval > 0 {
				time.Sleep(interval)
			}
			callCount++
			dr := destResult{Destination: dest.Address}
			if detail, err := amap.GetDrivingDetail(point, dest.Coordinate); err == nil {
				dr.Duration = formatDuration(detail.Duration)
				dr.Distance = formatDistance(detail.DistanceMeters)
				dr.Seconds = int64(detail.Duration.Seconds())
				dr.Meters = detail.DistanceMeters
			} else {
				dr.Error = err.Error()
			}
			item.Destinations[j] = dr
		}

		results[i] = item
	}

	var filePoints []filePointRecord
	for _, r := range results {
		var dests []fileDestRecord
		for j, d := range r.Destinations {
			if d.Error != "" {
				continue
			}
			dests = append(dests, fileDestRecord{
				Address: destinations[j].Address, Coordinate: destinations[j].Coordinate,
				Duration: d.Duration, Seconds: d.Seconds,
				Distance: d.Distance, Meters: d.Meters,
			})
		}
		if len(dests) > 0 {
			filePoints = append(filePoints, filePointRecord{
				Address: r.Address, Coordinate: r.Origin, Destinations: dests,
			})
		}
	}
	if len(filePoints) > 0 {
		appendResultsToFile(filePoints)
	}

	writeJSON(w, map[string]interface{}{"results": results})
}

// ---------- 结果文件 ----------

type fileRecord struct {
	Time   string            `json:"time"`
	Points []filePointRecord `json:"points"`
}

type filePointRecord struct {
	Address      string           `json:"address"`
	Coordinate   string           `json:"coordinate"`
	Destinations []fileDestRecord `json:"destinations"`
}

type fileDestRecord struct {
	Address    string `json:"address"`
	Coordinate string `json:"coordinate"`
	Duration   string `json:"duration"`
	Seconds    int64  `json:"seconds"`
	Distance   string `json:"distance"`
	Meters     int64  `json:"meters"`
}

const resultsFile = "results.json"

func appendResultsToFile(points []filePointRecord) {
	var records []json.RawMessage
	if data, err := os.ReadFile(resultsFile); err == nil && len(data) > 0 {
		json.Unmarshal(data, &records)
	}

	rec := fileRecord{
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		Points: points,
	}
	newJSON, err := json.Marshal(rec)
	if err != nil {
		log.Printf("序列化 JSON 失败: %v", err)
		return
	}
	records = append(records, json.RawMessage(newJSON))

	out, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		log.Printf("格式化 JSON 失败: %v", err)
		return
	}
	if err := os.WriteFile(resultsFile, out, 0644); err != nil {
		log.Printf("写入 %s 失败: %v", resultsFile, err)
	}
}

func handleResults(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(resultsFile)
	if err != nil || len(data) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte("[]"))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

// ---------- 工具 ----------

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%d小时%d分%d秒", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%d分%d秒", m, s)
	}
	return fmt.Sprintf("%d秒", s)
}

func formatDistance(meters int64) string {
	if meters < 1000 {
		return fmt.Sprintf("%d米", meters)
	}
	return fmt.Sprintf("%.1f公里", float64(meters)/1000)
}
