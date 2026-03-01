package amap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	geocodeURL    = "https://restapi.amap.com/v3/geocode/geo"
	regeoCodeURL  = "https://restapi.amap.com/v3/geocode/regeo"
)

// GeoCodeResult 地理编码结果
type GeoCodeResult struct {
	Country  interface{} `json:"country"`
	Province interface{} `json:"province"`
	City     interface{} `json:"city"`
	CityCode interface{} `json:"citycode"`
	District interface{} `json:"district"`
	Street   interface{} `json:"street"`
	Number   interface{} `json:"number"`
	AdCode   interface{} `json:"adcode"`
	Location string      `json:"location"`
	Level    string      `json:"level"`
}

func str(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func (g *GeoCodeResult) GetProvince() string { return str(g.Province) }
func (g *GeoCodeResult) GetCity() string     { return str(g.City) }
func (g *GeoCodeResult) GetDistrict() string { return str(g.District) }
func (g *GeoCodeResult) GetStreet() string   { return str(g.Street) }
func (g *GeoCodeResult) GetNumber() string   { return str(g.Number) }

type geocodeResponse struct {
	Status   string          `json:"status"`
	Info     string          `json:"info"`
	Count    string          `json:"count"`
	Geocodes []GeoCodeResult `json:"geocodes"`
}

// GeoCode 地理编码：将结构化地址转换为经纬度坐标。
// address 如 "北京市朝阳区阜通东大街6号"，city 可为空（全国搜索）。
func GeoCode(address, city string) ([]GeoCodeResult, error) {
	params := url.Values{}
	params.Set("key", AmapKey)
	params.Set("address", address)
	if city != "" {
		params.Set("city", city)
	}

	resp, err := http.Get(geocodeURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("请求地理编码 API 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result geocodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("API 返回错误: %s", result.Info)
	}

	if len(result.Geocodes) == 0 {
		return nil, fmt.Errorf("未找到匹配的地理编码结果")
	}

	return result.Geocodes, nil
}

// ReGeoResult 逆地理编码结果
type ReGeoResult struct {
	FormattedAddress string           `json:"formatted_address"`
	AddressComponent AddressComponent `json:"addressComponent"`
}

type AddressComponent struct {
	Country       string         `json:"country"`
	Province      string         `json:"province"`
	City          interface{}    `json:"city"`    // 直辖市时可能为空数组
	CityCode      string         `json:"citycode"`
	District      string         `json:"district"`
	AdCode        string         `json:"adcode"`
	Township      string         `json:"township"`
	TownCode      string         `json:"towncode"`
	Neighborhood  Neighborhood   `json:"neighborhood"`
	Building      Building       `json:"building"`
	StreetNumber  StreetNumber   `json:"streetNumber"`
	BusinessAreas []BusinessArea `json:"businessAreas"`
}

type Neighborhood struct {
	Name string      `json:"name"`
	Type interface{} `json:"type"`
}

type Building struct {
	Name string      `json:"name"`
	Type interface{} `json:"type"`
}

type StreetNumber struct {
	Street    string `json:"street"`
	Number    string `json:"number"`
	Location  string `json:"location"`
	Direction string `json:"direction"`
	Distance  string `json:"distance"`
}

type BusinessArea struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	ID       string `json:"id"`
}

type regeoResponse struct {
	Status    string      `json:"status"`
	Info      string      `json:"info"`
	ReGeoCode ReGeoResult `json:"regeocode"`
}

// CityName 安全地获取城市名称（直辖市时 city 字段为空数组）
func (ac *AddressComponent) CityName() string {
	if s, ok := ac.City.(string); ok {
		return s
	}
	return ""
}

// ReGeoCode 逆地理编码：将经纬度坐标转换为结构化地址。
// location 格式为 "经度,纬度"，例如 "116.481028,39.989643"。
func ReGeoCode(location string) (*ReGeoResult, error) {
	params := url.Values{}
	params.Set("key", AmapKey)
	params.Set("location", location)
	params.Set("extensions", "base")

	resp, err := http.Get(regeoCodeURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("请求逆地理编码 API 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result regeoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("API 返回错误: %s", result.Info)
	}

	if result.ReGeoCode.FormattedAddress == "" {
		return nil, fmt.Errorf("未找到匹配的逆地理编码结果")
	}

	return &result.ReGeoCode, nil
}
