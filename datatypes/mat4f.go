package datatypes

import (
	"fmt"
	"math"
)

type Mat4f struct {
	X Vec4f
	Y Vec4f
	Z Vec4f
	W Vec4f
}

func NewMat4f(x Vec4f, y Vec4f, z Vec4f, w Vec4f) Mat4f {
	return Mat4f{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}
}

func Scale(scale Vec3f) Mat4f {
	return NewMat4f(
		NewVec4f(scale.X, 0.0, 0.0, 0.0),
		NewVec4f(0.0, scale.Y, 0.0, 0.0),
		NewVec4f(0.0, 0.0, scale.Z, 0.0),
		NewVec4f(0.0, 0.0, 0.0, 1.0),
	)
}

func Translation(position Vec3f) Mat4f {
	return NewMat4f(
		NewVec4f(1.0, 0.0, 0.0, 0.0),
		NewVec4f(0.0, 1.0, 0.0, 0.0),
		NewVec4f(0.0, 0.0, 1.0, 0.0),
		NewVec4f(position.X, position.Y, -position.Z, 1.0),
	)
}

func RotationX(rotation float32) Mat4f {
	cos := float32(math.Cos(float64(rotation)))
	sin := float32(math.Sin(float64(rotation)))
	return NewMat4f(
		NewVec4f(1.0, 0.0, 0.0, 0.0),
		NewVec4f(0.0, cos, sin, 0.0),
		NewVec4f(0.0, -sin, cos, 0.0),
		NewVec4f(0.0, 0.0, 0.0, 1.0),
	)
}

func RotationY(rotation float32) Mat4f {
	cos := float32(math.Cos(float64(rotation)))
	sin := float32(math.Sin(float64(rotation)))
	return NewMat4f(
		NewVec4f(cos, 0.0, -sin, 0.0),
		NewVec4f(0.0, 1.0, 0.0, 0.0),
		NewVec4f(sin, 0.0, cos, 0.0),
		NewVec4f(0.0, 0.0, 0.0, 1.0),
	)
}

func RotationZ(rotation float32) Mat4f {
	cos := float32(math.Cos(float64(rotation)))
	sin := float32(math.Sin(float64(rotation)))
	return NewMat4f(
		NewVec4f(cos, sin, 0.0, 0.0),
		NewVec4f(-sin, cos, 0.0, 0.0),
		NewVec4f(0.0, 0.0, 1.0, 0.0),
		NewVec4f(0.0, 0.0, 0.0, 1.0),
	)
}

func Rotation(rotation Vec3f) Mat4f {
	return RotationX(rotation.X).Mul(RotationY(rotation.Y)).Mul(RotationZ(rotation.Z))
}

func Transformation(scale Vec3f, translation Vec3f, rotation Vec3f) Mat4f {
	return Scale(scale).Mul(Rotation(rotation)).Mul(Translation(translation))
}

func TransformationQuat(scale Vec3f, translation Vec3f, rotation Quaternion) Mat4f {
	return Scale(scale).Mul(rotation.ToMatrix()).Mul(Translation(translation))
}

func Orthographic(left float32, right float32, bottom float32, top float32, near float32, far float32) Mat4f {
	zero := float32(0.0)
	one := float32(1.0)

	c0r0 := 2.0 / (right - left)
	c0r1 := zero
	c0r2 := zero
	c0r3 := zero

	c1r0 := zero
	c1r1 := 2.0 / (top - bottom)
	c1r2 := zero
	c1r3 := zero

	c2r0 := zero
	c2r1 := zero
	c2r2 := -2.0 / (far - near)
	c2r3 := zero

	c3r0 := -(right + left) / (right - left)
	c3r1 := -(top + bottom) / (top - bottom)
	c3r2 := -(far + near) / (far - near)
	c3r3 := one

	return NewMat4f(
		NewVec4f(c0r0, c0r1, c0r2, c0r3),
		NewVec4f(c1r0, c1r1, c1r2, c1r3),
		NewVec4f(c2r0, c2r1, c2r2, c2r3),
		NewVec4f(c3r0, c3r1, c3r2, c3r3),
	)
}

func Perspective(fovY float32, aspect float32, near float32, far float32) (Mat4f, error) {
	if fovY <= 0.0 {
		return Mat4f{}, fmt.Errorf("the vertical field of view cannot be below zero, found: %v", fovY)
	}
	if near <= 0.0 {
		return Mat4f{}, fmt.Errorf("the near plane distance cannot be below zero, found: %v", near)
	}
	if far <= 0.0 {
		return Mat4f{}, fmt.Errorf("the far plane distance cannot be below zero, found: %v", far)
	}

	invLength := 1.0 / (near - far)
	f := 1.0 / float32(math.Tan(float64(0.5*fovY)))
	a := f / aspect
	b := (near + far) * invLength
	c := (2.0 * near * far) * invLength

	return NewMat4f(
		NewVec4f(a, 0.0, 0.0, 0.0),
		NewVec4f(0.0, f, 0.0, 0.0),
		NewVec4f(0.0, 0.0, b, -1.0),
		NewVec4f(0.0, 0.0, c, 0.0),
	), nil
}

func Identity() Mat4f {
	return NewMat4f(
		NewVec4f(1.0, 0.0, 0.0, 0.0),
		NewVec4f(0.0, 1.0, 0.0, 0.0),
		NewVec4f(0.0, 0.0, 1.0, 0.0),
		NewVec4f(0.0, 0.0, 0.0, 1.0),
	)
}

func (mat Mat4f) Mul(other Mat4f) Mat4f {
	return NewMat4f(
		NewVec4f(
			mat.X.X*other.X.X+mat.X.Y*other.Y.X+mat.X.Z*other.Z.X+mat.X.W*other.W.X,
			mat.X.X*other.X.Y+mat.X.Y*other.Y.Y+mat.X.Z*other.Z.Y+mat.X.W*other.W.Y,
			mat.X.X*other.X.Z+mat.X.Y*other.Y.Z+mat.X.Z*other.Z.Z+mat.X.W*other.W.Z,
			mat.X.X*other.X.W+mat.X.Y*other.Y.W+mat.X.Z*other.Z.W+mat.X.W*other.W.W,
		),
		NewVec4f(
			mat.Y.X*other.X.X+mat.Y.Y*other.Y.X+mat.Y.Z*other.Z.X+mat.Y.W*other.W.X,
			mat.Y.X*other.X.Y+mat.Y.Y*other.Y.Y+mat.Y.Z*other.Z.Y+mat.Y.W*other.W.Y,
			mat.Y.X*other.X.Z+mat.Y.Y*other.Y.Z+mat.Y.Z*other.Z.Z+mat.Y.W*other.W.Z,
			mat.Y.X*other.X.W+mat.Y.Y*other.Y.W+mat.Y.Z*other.Z.W+mat.Y.W*other.W.W,
		),
		NewVec4f(
			mat.Z.X*other.X.X+mat.Z.Y*other.Y.X+mat.Z.Z*other.Z.X+mat.Z.W*other.W.X,
			mat.Z.X*other.X.Y+mat.Z.Y*other.Y.Y+mat.Z.Z*other.Z.Y+mat.Z.W*other.W.Y,
			mat.Z.X*other.X.Z+mat.Z.Y*other.Y.Z+mat.Z.Z*other.Z.Z+mat.Z.W*other.W.Z,
			mat.Z.X*other.X.W+mat.Z.Y*other.Y.W+mat.Z.Z*other.Z.W+mat.Z.W*other.W.W,
		),
		NewVec4f(
			mat.W.X*other.X.X+mat.W.Y*other.Y.X+mat.W.Z*other.Z.X+mat.W.W*other.W.X,
			mat.W.X*other.X.Y+mat.W.Y*other.Y.Y+mat.W.Z*other.Z.Y+mat.W.W*other.W.Y,
			mat.W.X*other.X.Z+mat.W.Y*other.Y.Z+mat.W.Z*other.Z.Z+mat.W.W*other.W.Z,
			mat.W.X*other.X.W+mat.W.Y*other.Y.W+mat.W.Z*other.Z.W+mat.W.W*other.W.W,
		),
	)
}
