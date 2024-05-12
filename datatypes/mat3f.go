package datatypes

type Mat3f struct {
	X Vec3f
	Y Vec3f
	Z Vec3f
}

func NewMat3f(x Vec3f, y Vec3f, z Vec3f) Mat3f {
	return Mat3f{
		X: x,
		Y: y,
		Z: z,
	}
}

func Identity3() Mat3f {
	return NewMat3f(
		NewVec3f(1.0, 0.0, 0.0),
		NewVec3f(0.0, 1.0, 0.0),
		NewVec3f(0.0, 0.0, 1.0),
	)
}
