package value

import (
	"strconv"
)

// 长度
type Length float64

// 长度单位
const (
	Millimeter Length = 1                 // 毫米
	Meter             = 1000 * Millimeter // 米
	Kilometer         = 1000 * Meter      // 千米
)

func (d Length) String() string {
	return d.Format(Millimeter, 1)
}

// 格式化长度表示
// unit：单位
// precision：精度
func (d Length) Format(unit Length, precision int) string {
	switch unit {
	case Kilometer:
		return strconv.FormatFloat(float64(d/1000/1000), 'f', precision, 64) + "km"
	case Meter:
		return strconv.FormatFloat(float64(d/1000), 'f', precision, 64) + "m"
	default:
		return strconv.FormatFloat(float64(d), 'f', precision, 64) + "mm"
	}
}
