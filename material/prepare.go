package material

import "github.com/rajveermalviya/go-webgpu/wgpu"

func (m *Material) Prepare(renderPass *wgpu.RenderPassEncoder) {
	m.DiffuseTexture.SetBindGroup(renderPass)
}
