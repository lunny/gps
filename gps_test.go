package gps

import (
	"fmt"
	"testing"
)

func TestTransform(t *testing.T) {
	//国测局坐标(火星坐标,比如高德地图在用),百度坐标,wgs84坐标(谷歌国外以及绝大部分国外在线地图使用的坐标)
	//百度经纬度坐标转国测局坐标
	lat, lon := BD09ToGCJ02(39.915, 116.404)
	fmt.Println(lat, lon)

	//国测局坐标转百度经纬度坐标
	lat, lon = GCJ02ToBD09(39.915, 116.404)
	fmt.Println(lat, lon)

	//wgs84转国测局坐标
	lat, lon = WGS84ToGCJ02(39.915, 116.404)
	fmt.Println(lat, lon)

	//国测局坐标转wgs84坐标
	lat, lon = GCJ02ToWGS84(39.915, 116.404)
	fmt.Println(lat, lon)

	//result
	//bd09togcj02:   [ 116.39762729119315, 39.90865673957631 ]
	//gcj02tobd09:   [ 116.41036949371029, 39.92133699351021 ]
	//wgs84togcj02:  [ 116.41024449916938, 39.91640428150164 ]
	//gcj02towgs84:  [ 116.39775550083061, 39.91359571849836 ]
}
