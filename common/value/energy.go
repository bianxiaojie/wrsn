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
	return e.Format(Millijoule, 1)
}

// 格式化能量表示
// unit：单位
// precision：精度
func (e Energy) Format(unit Energy, precision int) string {
	switch unit {
	case Kilojoule:
		return strconv.FormatFloat(float64(e/1000/1000), 'f', precision, 64) + "kj"
	case Joule:
		return strconv.FormatFloat(float64(e/1000), 'f', precision, 64) + "j"
	default:
		return strconv.FormatFloat(float64(e), 'f', precision, 64) + "mj"
	}
}
