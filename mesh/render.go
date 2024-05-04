package mesh

import "github.com/rajveermalviya/go-webgpu/wgpu"

func (m *Mesh) DrawInstanced(instances uint32, renderPass *wgpu.RenderPassEncoder) {
	renderPass.SetVertexBuffer(0, m.VertexBuffer, 0, wgpu.WholeSize)
	renderPass.SetIndexBuffer(m.IndexBuffer, wgpu.IndexFormat_Uint32, 0, wgpu.WholeSize)
	renderPass.DrawIndexed(m.NumElements, instances, 0, 0, 0)
}

func (m *Mesh) Draw(renderPass *wgpu.RenderPassEncoder) {
	m.DrawInstanced(1, renderPass)
}
