package instance

import "porridgo/datatypes"

type Instance struct {
	Position datatypes.Vec3f
	Rotation datatypes.Quaternion
	Scale    datatypes.Vec3f
}
