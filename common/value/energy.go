package value

import (
	"strconv"
)

// 能量
type Energy float64

// 能量单位
const (
	Millijoule Energy = 1
	Joule             = 1000 * Millijoule
	Kilojoule         = 1000 * Joule
)

func (e Energy) String() string {
	return e.Format(Joule, 1)
}

// 格式化能量表示
// unit：单位
// precision：精度
func (e Energy) Format(unit Energy, precision int) string {
	switch unit {
	case Kilojoule:
		return strconv.FormatFloat(float64(e/1000/1000), 'f', precision, 64) + "kJ"
	case Joule:
		return strconv.FormatFloat(float64(e/1000), 'f', precision, 64) + "J"
	default:
		return strconv.FormatFloat(float64(e), 'f', precision, 64) + "mJ"
	}
}
