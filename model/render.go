package model

import (
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (m *Model) DrawInstanced(instances uint32, renderPass *wgpu.RenderPassEncoder) {
	for _, mesh := range m.Meshes {
		mtl := m.Materials[*mesh.Material]
		mesh.DrawInstanced(&mtl, instances, renderPass)
	}
}

func (m *Model) Draw(renderPass *wgpu.RenderPassEncoder) {
	m.DrawInstanced(1, renderPass)
}
