package datatypes

type Vec4f struct {
	X float32
	Y float32
	Z float32
	W float32
}

func NewVec4f(x float32, y float32, z float32, w float32) Vec4f {
	return Vec4f{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}
}
