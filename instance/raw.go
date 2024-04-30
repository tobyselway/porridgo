package instance

import "porridgo/datatypes"

type Raw struct {
	Model datatypes.Mat4f
}

func (inst Instance) ToRaw() Raw {
	return Raw{
		Model: datatypes.Transformation(inst.Scale, inst.Position, inst.Rotation),
	}
}
