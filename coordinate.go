package amap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const convertURL = "https://restapi.amap.com/v3/assistant/coordinate/convert"

// CoordSys 原始坐标系类型
type CoordSys string

const (
	CoordGPS     CoordSys = "gps"
	CoordMapbar  CoordSys = "mapbar"
	CoordBaidu   CoordSys = "baidu"
	CoordAutoNavi CoordSys = "autonavi"
)

type convertResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Locations string `json:"locations"`
}

// ConvertCoordinate 将其他坐标系的坐标转换为高德坐标（GCJ-02）。
// locations 为坐标列表，格式 "经度,纬度"；多个坐标用 "|" 分隔，最多40对。
// coordsys 为原始坐标系（gps / mapbar / baidu / autonavi）。
// 返回转换后的坐标切片，每个元素格式为 "经度,纬度"。
func ConvertCoordinate(locations string, coordsys CoordSys) ([]string, error) {
	params := url.Values{}
	params.Set("key", AmapKey)
	params.Set("locations", locations)
	params.Set("coordsys", string(coordsys))

	resp, err := http.Get(convertURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("请求坐标转换 API 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result convertResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("API 返回错误: %s", result.Info)
	}

	if result.Locations == "" {
		return nil, fmt.Errorf("未返回转换结果")
	}

	return strings.Split(result.Locations, ";"), nil
}
