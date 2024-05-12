package pipeline

import (
	"porridgo/label"
	"porridgo/texture"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (p *Pipeline) Setup(device *wgpu.Device, swapChainConfig *wgpu.SwapChainDescriptor) error {
	err := p.VertexShader.Setup(device)
	if err != nil {
		return err
	}
	err = p.FragmentShader.Setup(device)
	if err != nil {
		return err
	}
	err = p.setupPipelineLayout(device)
	if err != nil {
		return err
	}
	err = p.createPipeline(device, swapChainConfig)
	if err != nil {
		return err
	}
	return nil
}

func (p *Pipeline) createPipeline(device *wgpu.Device, swapChainConfig *wgpu.SwapChainDescriptor) error {
	var err error
	p.renderPipeline, err = device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label:  label.Label(p, "Pipeline"),
		Layout: p.layout,
		Vertex: wgpu.VertexState{
			Module:     p.VertexShader.Module,
			EntryPoint: "vs_main",
			Buffers:    p.VertexBufferLayouts,
		},
		Primitive: wgpu.PrimitiveState{
			Topology:         wgpu.PrimitiveTopology_TriangleList,
			StripIndexFormat: wgpu.IndexFormat_Undefined,
			FrontFace:        wgpu.FrontFace_CCW,
			CullMode:         wgpu.CullMode_None,
		},
		Multisample: wgpu.MultisampleState{
			Count:                  1,
			Mask:                   0xFFFFFFFF,
			AlphaToCoverageEnabled: false,
		},
		Fragment: &wgpu.FragmentState{
			Module:     p.FragmentShader.Module,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{
				{
					Format:    swapChainConfig.Format,
					Blend:     &wgpu.BlendState_Replace,
					WriteMask: wgpu.ColorWriteMask_All,
				},
			},
		},
		DepthStencil: &wgpu.DepthStencilState{
			Format:            texture.DEPTH_FORMAT,
			DepthWriteEnabled: true,
			DepthCompare:      wgpu.CompareFunction_Less,
			StencilFront: wgpu.StencilFaceState{
				Compare: wgpu.CompareFunction_Never,
			},
			StencilBack: wgpu.StencilFaceState{
				Compare: wgpu.CompareFunction_Never,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
