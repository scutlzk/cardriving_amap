package main

import (
	"fmt"
	"log"

	"amap"
)

func main() {
	origin := "116.481028,39.989643"
	destination := "116.434446,39.90816"

	// 1. 驾车路线规划
	fmt.Println("========== 驾车路线规划 ==========")
	duration, err := amap.DrivingDuration(origin, destination)
	if err != nil {
		log.Fatal(err)
	}
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60
	fmt.Printf("起点: %s\n终点: %s\n驾车最短用时: %d分%d秒\n", origin, destination, minutes, seconds)

	// 2. 逆地理编码：经纬度 → 地址
	fmt.Println("\n========== 逆地理编码 ==========")
	fmt.Printf("坐标: %s\n", origin)
	regeo, err := amap.ReGeoCode(origin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("地址: %s\n", regeo.FormattedAddress)
	ac := regeo.AddressComponent
	fmt.Printf("省份: %s  城市: %s  区县: %s  乡镇/街道: %s\n",
		ac.Province, ac.CityName(), ac.District, ac.Township)

	// 3. 地理编码：地址 → 经纬度（用逆地理编码得到的地址做正向验证）
	fmt.Println("\n========== 地理编码 ==========")
	address := regeo.FormattedAddress
	fmt.Printf("地址: %s\n", address)
	geocodes, err := amap.GeoCode(address, "")
	if err != nil {
		log.Fatal(err)
	}
	for i, g := range geocodes {
		fmt.Printf("结果%d: 坐标=%s  级别=%s\n", i+1, g.Location, g.Level)
	}

	// 4. 行政区域查询：东莞市虎门区，subdistrict=2 查看下两级
	fmt.Println("\n========== 行政区域查询 ==========")
	keyword := "虎门"
	fmt.Printf("关键词: %s\n", keyword)
	districts, err := amap.QueryDistrict(keyword, 1)
	if err != nil {
		log.Fatal(err)
	}
	var printDistrict func(d amap.District, indent string)
	printDistrict = func(d amap.District, indent string) {
		fmt.Printf("%s%s  (adcode=%s, center=%s, level=%s)\n",
			indent, d.Name, d.AdCode, d.Center, d.Level)
		for _, sub := range d.Districts {
			printDistrict(sub, indent+"  ")
		}
	}
	for _, d := range districts {
		printDistrict(d, "")
	}

	// 5. 坐标转换：GPS坐标 → 高德坐标（GCJ-02）
	fmt.Println("\n========== 坐标转换 ==========")
	gpsCoord := "116.481499,39.990475"
	fmt.Printf("GPS 原始坐标: %s\n", gpsCoord)
	converted, err := amap.ConvertCoordinate(gpsCoord, amap.CoordGPS)
	if err != nil {
		log.Fatal(err)
	}
	for i, c := range converted {
		fmt.Printf("高德坐标%d: %s\n", i+1, c)
	}
}
