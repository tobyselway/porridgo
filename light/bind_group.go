package light

import (
	"fmt"
	"porridgo/label"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

var BindGroupLayout *wgpu.BindGroupLayout = nil

func SetupBindGroupLayout(device *wgpu.Device) error {
	if BindGroupLayout != nil {
		return fmt.Errorf("bind group layout already created")
	}
	var err error = nil
	BindGroupLayout, err = device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Label: "Light Bind Group Layout",
		Entries: []wgpu.BindGroupLayoutEntry{
			{
				Binding:    0,
				Visibility: wgpu.ShaderStage_Vertex | wgpu.ShaderStage_Fragment,
				Buffer: wgpu.BufferBindingLayout{
					Type:             wgpu.BufferBindingType_Uniform,
					HasDynamicOffset: false,
					MinBindingSize:   0,
				},
			},
		},
	})
	return err
}

func (l *Light) CreateBindGroup(device *wgpu.Device) error {
	var err error
	l.bindGroup, err = device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  label.Label(l, "Bind Group"),
		Layout: BindGroupLayout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding: 0,
				Buffer:  l.buffer,
				Size:    l.buffer.GetSize(),
			},
		},
	})
	return err
}

func CleanupBindGroupLayout() {
	BindGroupLayout.Release()
	BindGroupLayout = nil
}

func (l *Light) SetBindGroup(groupIndex uint32, renderPass *wgpu.RenderPassEncoder) {
	renderPass.SetBindGroup(groupIndex, l.bindGroup, nil)
}
