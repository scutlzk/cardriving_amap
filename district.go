package amap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const districtURL = "https://restapi.amap.com/v3/config/district"

// District 行政区信息
type District struct {
	CityCode  string     `json:"citycode"`
	AdCode    string     `json:"adcode"`
	Name      string     `json:"name"`
	Center    string     `json:"center"`
	Level     string     `json:"level"`
	Districts []District `json:"districts"`
}

type districtResponse struct {
	Status   string     `json:"status"`
	Info     string     `json:"info"`
	InfoCode string     `json:"infocode"`
	Districts []District `json:"districts"`
}

// QueryDistrict 行政区域查询。
// keywords 支持行政区名称、citycode、adcode；subdistrict 为返回的下级行政区层级数（0~3）。
func QueryDistrict(keywords string, subdistrict int) ([]District, error) {
	params := url.Values{}
	params.Set("key", AmapKey)
	params.Set("keywords", keywords)
	params.Set("subdistrict", fmt.Sprintf("%d", subdistrict))

	resp, err := http.Get(districtURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("请求行政区域查询 API 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result districtResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("API 返回错误: %s (infocode=%s)", result.Info, result.InfoCode)
	}

	if len(result.Districts) == 0 {
		return nil, fmt.Errorf("未找到匹配的行政区域")
	}

	return result.Districts, nil
}
