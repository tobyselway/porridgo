package pipeline

import (
	"fmt"
	"porridgo/label"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (p *Pipeline) setupPipelineLayout(device *wgpu.Device) error {
	if p.layout != nil {
		return fmt.Errorf("pipeline layout already created")
	}
	var err error = nil
	p.layout, err = device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label:              label.Label(p, "Pipeline Layout"),
		BindGroupLayouts:   p.BindGroupLayouts,
		PushConstantRanges: []wgpu.PushConstantRange{},
	})
	return err
}

func (p *Pipeline) cleanupPipelineLayout() {
	if p.layout != nil {
		defer p.layout.Release()
		p.layout = nil
	}
}
