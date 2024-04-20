package datatypes

type Vec2f struct {
	X float32
	Y float32
}

func NewVec2f(x float32, y float32) Vec2f {
	return Vec2f{
		X: x,
		Y: y,
	}
}
