package pipeline

import "github.com/rajveermalviya/go-webgpu/wgpu"

func (p *Pipeline) Prepare(renderPass *wgpu.RenderPassEncoder) {
	renderPass.SetPipeline(p.renderPipeline)
}
