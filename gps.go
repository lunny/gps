/**
 * 各地图API坐标系统比较与转换;
 * WGS84坐标系：即地球坐标系，国际上通用的坐标系。设备一般包含GPS芯片或者北斗芯片获取的经纬度为WGS84地理坐标系,
 * 谷歌地图采用的是WGS84地理坐标系（中国范围除外）;
 * GCJ02坐标系：即火星坐标系，是由中国国家测绘局制订的地理信息系统的坐标系统。由WGS84坐标系经加密后的坐标系。
 * 谷歌中国地图和搜搜中国地图采用的是GCJ02地理坐标系; BD09坐标系：即百度坐标系，GCJ02坐标系经加密后的坐标系;
 * 搜狗坐标系、图吧坐标系等，估计也是在GCJ02基础上加密而成的。 chenhua
 */

package gps

import (
	"math"
)

// some constant variables
const (
	XPI float64 = 3.14159265358979324 * 3000.0 / 180.0
	PI  float64 = 3.1415926535897932384626
	A   float64 = 6378245.0
	EE  float64 = 0.00669342162296594323
)

// OutOfChina returns true if an gps84 position is not in China
func OutOfChina(lat, lon float64) bool {
	return ((lon < 72.004 || lon > 137.8347) && (lat < 0.8293 || lat > 55.8271))
}

// WGS84ToGCJ02 transform World Geodetic System ==> Mars Geodetic System
func WGS84ToGCJ02(lat, lon float64) (float64, float64) {
	if OutOfChina(lat, lon) {
		return lat, lon
	}

	dLat := transformLat(lon-105.0, lat-35.0)
	dLon := transformLon(lon-105.0, lat-35.0)
	radLat := lat / 180.0 * PI
	magic := sin(radLat)
	magic = 1 - EE*magic*magic
	sqrtMagic := sqrt(magic)
	dLat = (dLat * 180.0) / ((A * (1 - EE)) / (magic * sqrtMagic) * PI)
	dLon = (dLon * 180.0) / (A / sqrtMagic * cos(radLat) * PI)
	mgLat := lat + dLat
	mgLon := lon + dLon
	return mgLat, mgLon
}

func WGS84ToBD09(lat, lon float64) (float64, float64) {
	return GCJ02ToBD09(WGS84ToGCJ02(lat, lon))
}

// GCJ02ToWGS84 火星坐标系 (GCJ02) to GPS84
func GCJ02ToWGS84(lat, lon float64) (float64, float64) {
	tLat, tLon := transform(lat, lon)
	lontitude := lon*2 - tLon
	latitude := lat*2 - tLat
	return latitude, lontitude
}

// GCJ02ToBD09 火星坐标系 (GCJ-02) 与百度坐标系 (BD-09) 的转换算法 将 GCJ-02 坐标转换成 BD-09 坐标
func GCJ02ToBD09(lat, lon float64) (float64, float64) {
	z := sqrt(lat*lat+lon*lon) + 0.00002*sin(lat*XPI)
	theta := atan2(lat, lon) + 0.000003*cos(lon*XPI)
	bdLon := z*cos(theta) + 0.0065
	bdLat := z*sin(theta) + 0.006
	return bdLat, bdLon
}

// BD09ToGCJ02 火星坐标系 (GCJ-02) 与百度坐标系 (BD-09) 的转换算法
func BD09ToGCJ02(bdLat, bdLon float64) (float64, float64) {
	x := bdLon - 0.0065
	y := bdLat - 0.006
	z := sqrt(x*x+y*y) - 0.00002*sin(y*XPI)
	theta := atan2(y, x) - 0.000003*cos(x*XPI)
	ggLon := z * cos(theta)
	ggLat := z * sin(theta)
	return ggLat, ggLon
}

// BD09ToGPS84 百度坐标系 转换为 国际坐标系
func BD09ToGPS84(bdLat, bdLon float64) (float64, float64) {
	gcj02Lat, gcj02Lon := BD09ToGCJ02(bdLat, bdLon)
	return GCJ02ToWGS84(gcj02Lat, gcj02Lon)
}

func transform(lat, lon float64) (float64, float64) {
	if OutOfChina(lat, lon) {
		return lat, lon
	}

	dLat := transformLat(lon-105.0, lat-35.0)
	dLon := transformLon(lon-105.0, lat-35.0)
	radLat := lat / 180.0 * PI
	magic := sin(radLat)
	magic = 1 - EE*magic*magic
	sqrtMagic := sqrt(magic)
	dLat = (dLat * 180.0) / ((A * (1 - EE)) / (magic * sqrtMagic) * PI)
	dLon = (dLon * 180.0) / (A / sqrtMagic * cos(radLat) * PI)
	mgLat := lat + dLat
	mgLon := lon + dLon
	return mgLat, mgLon
}

func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*sqrt(abs(x))
	ret += (20.0*sin(6.0*x*PI) + 20.0*sin(2.0*x*PI)) * 2.0 / 3.0
	ret += (20.0*sin(y*PI) + 40.0*sin(y/3.0*PI)) * 2.0 / 3.0
	ret += (160.0*sin(y/12.0*PI) + 320*sin(y*PI/30.0)) * 2.0 / 3.0
	return ret
}

func transformLon(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*sqrt(abs(x))
	ret += (20.0*sin(6.0*x*PI) + 20.0*sin(2.0*x*PI)) * 2.0 / 3.0
	ret += (20.0*sin(x*PI) + 40.0*sin(x/3.0*PI)) * 2.0 / 3.0
	ret += (150.0*sin(x/12.0*PI) + 300.0*sin(x/30.0*PI)) * 2.0 / 3.0
	return ret
}

func atan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

func cos(x float64) float64 {
	return math.Cos(x)
}

func sin(x float64) float64 {
	return math.Sin(x)
}

func sqrt(in float64) float64 {
	return math.Sqrt(in)
}

func abs(in float64) float64 {
	return math.Abs(in)
}
