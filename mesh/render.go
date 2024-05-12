package mesh

import (
	"porridgo/material"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (m *Mesh) DrawInstanced(mtl *material.Material, instances uint32, renderPass *wgpu.RenderPassEncoder) {
	renderPass.SetVertexBuffer(0, m.VertexBuffer, 0, wgpu.WholeSize)
	renderPass.SetIndexBuffer(m.IndexBuffer, wgpu.IndexFormat_Uint32, 0, wgpu.WholeSize)
	if mtl != nil {
		mtl.Prepare(renderPass)
	}
	renderPass.DrawIndexed(m.NumElements, instances, 0, 0, 0)
}

func (m *Mesh) Draw(mtl *material.Material, renderPass *wgpu.RenderPassEncoder) {
	m.DrawInstanced(mtl, 1, renderPass)
}
