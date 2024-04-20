package datatypes

type Vec3f struct {
	X float32
	Y float32
	Z float32
}

func NewVec3f(x float32, y float32, z float32) Vec3f {
	return Vec3f{
		X: x,
		Y: y,
		Z: z,
	}
}
