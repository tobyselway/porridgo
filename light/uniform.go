package light

import (
	"porridgo/datatypes"
)

type Uniform struct {
	Position  datatypes.Vec3f
	_padding  float32
	Color     datatypes.Vec3f
	_padding2 float32
}

func (l *Light) SetupUniform() {
	l.uniform = Uniform{
		Position: datatypes.NewVec3f(2.0, 2.0, 2.0),
		Color:    datatypes.NewVec3f(1.0, 1.0, 1.0),
	}
}

func (l *Light) UpdateUniform() error {
	l.uniform.Position = l.Position
	l.uniform.Color = l.Color
	return nil
}
