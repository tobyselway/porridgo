package instance

import "porridgo/datatypes"

type Instance struct {
	Position datatypes.Vec3f
	// Rotation datatypes.Quaternion // TODO
	Rotation datatypes.Vec3f
	Scale    datatypes.Vec3f
}
