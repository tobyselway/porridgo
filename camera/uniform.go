package camera

import (
	"fmt"
	"porridgo/datatypes"
)

type Uniform struct {
	Position   datatypes.Vec4f
	View       datatypes.Mat4f
	Projection datatypes.Mat4f
}

func (c *Camera) SetupUniform() {
	c.uniform = Uniform{
		Position:   datatypes.NewVec4f(0, 0, 0, 0),
		View:       datatypes.Identity(),
		Projection: datatypes.Identity(),
	}
}

func (c *Camera) UpdateUniform() error {
	c.uniform.Position = c.Position.ToVec4f()
	c.uniform.View = c.BuildViewMatrix()
	projMat, err := c.BuildProjectionMatrix()
	if err != nil {
		return fmt.Errorf("updating camera uniform: %w", err)
	}
	c.uniform.Projection = projMat
	return nil
}
