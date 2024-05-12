package datatypes

import (
	"fmt"
	"math"
)

type Vec3f struct {
	X float32
	Y float32
	Z float32
}

func (v Vec3f) String() string {
	return fmt.Sprintf("[%v, %v, %v]", v.X, v.Y, v.Z)
}

func NewVec3f(x float32, y float32, z float32) Vec3f {
	return Vec3f{
		X: x,
		Y: y,
		Z: z,
	}
}

func (vec Vec3f) Add(other Vec3f) Vec3f {
	return NewVec3f(
		vec.X+other.X,
		vec.Y+other.Y,
		vec.Z+other.Z,
	)
}

func (vec Vec3f) Sub(other Vec3f) Vec3f {
	return NewVec3f(
		vec.X-other.X,
		vec.Y-other.Y,
		vec.Z-other.Z,
	)
}

func (vec Vec3f) MulElementWise(other Vec3f) Vec3f {
	return NewVec3f(
		vec.X*other.X,
		vec.Y*other.Y,
		vec.Z*other.Z,
	)
}

func (vec Vec3f) MulScalar(other float32) Vec3f {
	return NewVec3f(
		vec.X*other,
		vec.Y*other,
		vec.Z*other,
	)
}

func (vec Vec3f) Sum() float32 {
	return vec.X + vec.Y + vec.Z
}

func (vec Vec3f) Dot(other Vec3f) float32 {
	return vec.MulElementWise(other).Sum()
}

func (vec Vec3f) Cross(other Vec3f) Vec3f {
	return NewVec3f(
		(vec.Y*other.Z)-(vec.Z*other.Y),
		(vec.Z*other.X)-(vec.X*other.Z),
		(vec.X*other.Y)-(vec.Y*other.X),
	)
}

func (vec Vec3f) magnitude2() float32 {
	return vec.Dot(vec)
}

func (vec Vec3f) Magnitude() float32 {
	return float32(math.Sqrt(float64(vec.magnitude2())))
}

func (vec Vec3f) NormalizeTo(magnitude float32) Vec3f {
	if vec.Magnitude() == 0 {
		return NewVec3f(0.0, 0.0, 0.0)
	}
	return vec.MulScalar(magnitude / vec.Magnitude())
}

func (vec Vec3f) Normalize() Vec3f {
	return vec.NormalizeTo(1.0)
}

func (vec Vec3f) IsZero() bool {
	return vec.X == 0 && vec.Y == 0 && vec.Z == 0
}

func (vec Vec3f) ToVec4f() Vec4f {
	return NewVec4f(
		vec.X,
		vec.Y,
		vec.Z,
		1,
	)
}
