package amap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const drivingURL = "https://restapi.amap.com/v5/direction/driving"

// AmapKey 高德 Web 服务 API Key，可在程序启动时从配置覆盖。
var AmapKey = "09160f86bd4f4e3a8f9cb7797d4f5c7e"

type drivingResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Count    string `json:"count"`
	Route    struct {
		Origin      string `json:"origin"`
		Destination string `json:"destination"`
		Paths       []struct {
			Distance string `json:"distance"`
			Cost     struct {
				Duration      string `json:"duration"`
				Tolls         string `json:"tolls"`
				TollDistance   string `json:"toll_distance"`
				TrafficLights string `json:"traffic_lights"`
			} `json:"cost"`
		} `json:"paths"`
	} `json:"route"`
}

// DrivingDetail 驾车路线详情
type DrivingDetail struct {
	Duration       time.Duration
	DistanceMeters int64
}

// GetDrivingDetail 计算从 origin 到 destination 驾车最优路线的用时与距离。
// origin / destination 格式为 "经度,纬度"。
func GetDrivingDetail(origin, destination string) (*DrivingDetail, error) {
	params := url.Values{}
	params.Set("key", AmapKey)
	params.Set("origin", origin)
	params.Set("destination", destination)
	params.Set("strategy", "0")
	params.Set("show_fields", "cost")

	resp, err := http.Get(drivingURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("请求高德 API 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result drivingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("API 返回错误: %s (infocode=%s)", result.Info, result.Infocode)
	}

	if len(result.Route.Paths) == 0 {
		return nil, fmt.Errorf("未找到可用路线")
	}

	var best *DrivingDetail
	for _, p := range result.Route.Paths {
		sec, err := strconv.ParseInt(p.Cost.Duration, 10, 64)
		if err != nil {
			continue
		}
		dist, _ := strconv.ParseInt(p.Distance, 10, 64)
		if best == nil || sec < int64(best.Duration.Seconds()) {
			best = &DrivingDetail{
				Duration:       time.Duration(sec) * time.Second,
				DistanceMeters: dist,
			}
		}
	}

	if best == nil {
		return nil, fmt.Errorf("无法解析路线耗时")
	}

	return best, nil
}

// DrivingDuration 计算从 origin 到 destination 驾车最短用时。
func DrivingDuration(origin, destination string) (time.Duration, error) {
	detail, err := GetDrivingDetail(origin, destination)
	if err != nil {
		return 0, err
	}
	return detail.Duration, nil
}
