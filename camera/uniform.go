package camera

import (
	"fmt"
	"porridgo/datatypes"
)

type Uniform struct {
	View       datatypes.Mat4f
	Projection datatypes.Mat4f
}

func NewUniform() Uniform {
	return Uniform{
		View:       datatypes.Identity(),
		Projection: datatypes.Identity(),
	}
}

func (u *Uniform) Update(camera *Camera) error {
	u.View = camera.BuildViewMatrix()
	projMat, err := camera.BuildProjectionMatrix()
	if err != nil {
		return fmt.Errorf("updating camera uniform: %w", err)
	}
	u.Projection = projMat
	return nil
}
