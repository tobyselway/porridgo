package camera

import (
	"fmt"
	"porridgo/datatypes"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Camera struct {
	Position  datatypes.Vec3f
	Pitch     float32
	Yaw       float32
	Aspect    float32
	FovY      float32
	ZNear     float32
	ZFar      float32
	bindGroup *wgpu.BindGroup
	uniform   Uniform
	buffer    *wgpu.Buffer
}

func (c *Camera) Cleanup() {
	defer c.bindGroup.Release()
	defer c.buffer.Release()
}

var OPENGL_TO_WGPU_MATRIX datatypes.Mat4f = datatypes.NewMat4f(
	datatypes.NewVec4f(1.0, 0.0, 0.0, 0.0),
	datatypes.NewVec4f(0.0, 1.0, 0.0, 0.0),
	datatypes.NewVec4f(0.0, 0.0, 0.5, 0.5),
	datatypes.NewVec4f(0.0, 0.0, 0.0, 1.0),
)

func (c Camera) BuildViewMatrix() datatypes.Mat4f {
	return datatypes.Euler(datatypes.NewVec3f(c.Pitch, c.Yaw, 0.0)).ToMatrix().
		Mul(datatypes.Translation(c.Position.MulScalar(-1.0)))
}

func (c Camera) BuildProjectionMatrix() (datatypes.Mat4f, error) {
	// return datatypes.Orthographic(-1.0, 1.0, -1.0, 1.0, -5.0, 5.0), nil
	mat, err := datatypes.Perspective(c.FovY, c.Aspect, c.ZNear, c.ZFar)
	if err != nil {
		return datatypes.Mat4f{}, fmt.Errorf("building projection matrix: %w", err)
	}
	return mat, nil
}
