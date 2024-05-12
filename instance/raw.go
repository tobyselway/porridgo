package instance

import "porridgo/datatypes"

type Raw struct {
	Model  datatypes.Mat4f
	Normal datatypes.Mat3f
}

func (inst Instance) ToRaw() Raw {
	return Raw{
		Model:  datatypes.TransformationQuat(inst.Scale, inst.Position, inst.Rotation),
		Normal: inst.Rotation.ToMatrix().ToMat3f(),
	}
}
