package pipeline

import (
	"porridgo/model"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (p *Pipeline) DrawInstanced(m *model.Model, instances uint32, renderPass *wgpu.RenderPassEncoder) {
	m.DrawInstanced(instances, renderPass)
}

func (p *Pipeline) Draw(m *model.Model, renderPass *wgpu.RenderPassEncoder) {
	p.DrawInstanced(m, 1, renderPass)
}
