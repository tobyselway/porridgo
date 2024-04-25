package camera

import (
	"fmt"
	"porridgo/datatypes"
)

type Uniform struct {
	View       datatypes.Mat4f
	Projection datatypes.Mat4f
}

func (c *Camera) SetupUniform() {
	c.uniform = Uniform{
		View:       datatypes.Identity(),
		Projection: datatypes.Identity(),
	}
}

func (c *Camera) UpdateUniform() error {
	c.uniform.View = c.BuildViewMatrix()
	projMat, err := c.BuildProjectionMatrix()
	if err != nil {
		return fmt.Errorf("updating camera uniform: %w", err)
	}
	c.uniform.Projection = projMat
	return nil
}
