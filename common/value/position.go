package value

import (
	"fmt"
	"math"
)

// 位置
type Position struct {
	x Length // x坐标
	y Length // y坐标
}

func MakePosition(x Length, y Length) Position {
	return Position{
		x: x,
		y: y,
	}
}

func (p Position) X() Length {
	return p.x
}

func (p Position) Y() Length {
	return p.y
}

// 计算位置p与位置target的极角
func (p Position) PolarAngleTo(target Position) PolarAngle {
	xDiff, yDiff := float64(target.x-p.x), float64(target.y-p.y)
	acos := math.Acos(xDiff / math.Hypot(xDiff, yDiff))

	var pa PolarAngle
	if yDiff > 0 {
		pa = PolarAngle(acos)
	} else {
		pa = PolarAngle(2*math.Pi - acos)
	}
	return pa.Normalize()
}

// 计算位置p到位置target的距离
func (p Position) DistanceTo(target Position) Length {
	xDiff, yDiff := float64(target.x-p.x), float64(target.y-p.y)
	return Length(math.Hypot(xDiff, yDiff))
}

// 计算polarAngle方向length距离处位置坐标
func (p Position) PositionTo(polarAngle PolarAngle, length Length) Position {
	cos, sin := math.Cos(float64(polarAngle)), math.Sin(float64(polarAngle))
	x := p.x + length*Length(cos)
	y := p.y + length*Length(sin)
	return MakePosition(x, y)
}

func (p Position) String() string {
	return fmt.Sprintf("(%v, %v)", p.x, p.y)
}
