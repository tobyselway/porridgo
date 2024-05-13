package datatypes

import "fmt"

type Vec2f struct {
	X float32
	Y float32
}

func (v Vec2f) String() string {
	return fmt.Sprintf("[%v, %v]", v.X, v.Y)
}

func NewVec2f(x float32, y float32) Vec2f {
	return Vec2f{
		X: x,
		Y: y,
	}
}

func (vec Vec2f) Sub(other Vec2f) Vec2f {
	return NewVec2f(
		vec.X-other.X,
		vec.Y-other.Y,
	)
}
