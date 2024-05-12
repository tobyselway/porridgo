package pipeline

import (
	"porridgo/shader"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Pipeline struct {
	Name                string
	VertexShader        *shader.Shader
	FragmentShader      *shader.Shader
	BindGroupLayouts    []*wgpu.BindGroupLayout
	VertexBufferLayouts []wgpu.VertexBufferLayout
	renderPipeline      *wgpu.RenderPipeline
	layout              *wgpu.PipelineLayout
}

func (p *Pipeline) Cleanup() {
	defer p.cleanupPipelineLayout()
	if p.VertexShader != nil {
		defer p.VertexShader.Cleanup()
		p.VertexShader = nil
	}
	if p.FragmentShader != nil {
		defer p.FragmentShader.Cleanup()
		p.FragmentShader = nil
	}
	if p.renderPipeline != nil {
		defer p.renderPipeline.Release()
		p.renderPipeline = nil
	}
}
