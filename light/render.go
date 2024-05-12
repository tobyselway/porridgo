package light

import (
	"porridgo/model"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (l *Light) DrawInstanced(m *model.Model, instances uint32, renderPass *wgpu.RenderPassEncoder) {
	for _, mesh := range m.Meshes {
		mesh.DrawInstanced(nil, instances, renderPass)
	}
}

func (l *Light) Draw(m *model.Model, renderPass *wgpu.RenderPassEncoder) {
	l.DrawInstanced(m, 1, renderPass)
}
