package datatypes

import "math"

type Quaternion struct {
	Vec4f
}

func QuaternionFromVec4f(v Vec4f) Quaternion {
	return Quaternion{
		v,
	}
}

func NewQuaternion(x float32, y float32, z float32, w float32) Quaternion {
	return Quaternion{
		Vec4f{
			X: x,
			Y: y,
			Z: z,
			W: w,
		},
	}
}

func AngleAxis(axis Vec3f, angle float32) Quaternion {
	s := float32(math.Sin(float64(angle / 2.0)))
	return NewQuaternion(
		axis.X*s,
		axis.Y*s,
		axis.Z*s,
		float32(math.Cos(float64(angle/2.0))),
	)
}

func Euler(rotation Vec3f) Quaternion {
	c1 := float32(math.Cos(float64(rotation.Y / 2.0)))
	s1 := float32(math.Sin(float64(rotation.Y / 2.0)))
	c2 := float32(math.Cos(float64(rotation.Z / 2.0)))
	s2 := float32(math.Sin(float64(rotation.Z / 2.0)))
	c3 := float32(math.Cos(float64(rotation.X / 2.0)))
	s3 := float32(math.Sin(float64(rotation.X / 2.0)))
	c1c2 := c1 * c2
	s1s2 := s1 * s2
	return NewQuaternion(
		c1c2*s3+s1s2*c3,
		s1*c2*c3+c1*s2*s3,
		c1*s2*c3-s1*c2*s3,
		c1c2*c3-s1s2*s3,
	)
}

func (q Quaternion) ToMatrix() Mat4f {
	return NewMat4f(
		NewVec4f(
			1.0-2.0*(q.Y*q.Y)-2.0*(q.Z*q.Z),
			2.0*q.X*q.Y-2.0*q.W*q.Z,
			2.0*q.X*q.Z+2.0*q.W*q.Y,
			0.0,
		),
		NewVec4f(
			2.0*q.X*q.Y+2.0*q.W*q.Z,
			1.0-2.0*(q.X*q.X)-2.0*(q.Z*q.Z),
			2.0*q.Y*q.Z-2.0*q.W*q.X,
			0.0,
		),
		NewVec4f(
			2.0*q.X*q.Z-2.0*q.W*q.Y,
			2.0*q.Y*q.Z+2.0*q.W*q.X,
			1.0-2.0*(q.X*q.X)-2.0*(q.Y*q.Y),
			0.0,
		),
		NewVec4f(0.0, 0.0, 0.0, 1.0),
	)
}
