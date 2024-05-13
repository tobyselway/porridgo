package light

import (
	"porridgo/datatypes"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Light struct {
	Position  datatypes.Vec3f
	Color     datatypes.Vec3f
	bindGroup *wgpu.BindGroup
	uniform   Uniform
	buffer    *wgpu.Buffer
}

func (l *Light) Cleanup() {
	if l.bindGroup != nil {
		defer l.bindGroup.Release()
	}
	if l.buffer != nil {
		defer l.buffer.Release()
	}
}
