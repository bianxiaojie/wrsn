package value

import (
	"math"
	"strconv"
)

// 极角
type PolarAngle float64

func (pa PolarAngle) Normalize() PolarAngle {
	return PolarAngle(math.Mod(float64(pa), 2*math.Pi))
}

func (pa PolarAngle) String() string {
	return pa.Format(1)
}

// 格式化极角表示
// precision：精度
func (pa PolarAngle) Format(precision int) string {
	angle := float64(pa.Normalize() * 360 / (2 * math.Pi))
	return strconv.FormatFloat(angle, 'f', precision, 64) + "°"
}
