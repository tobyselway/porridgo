package camera

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
		Label: "Camera Bind Group Layout",
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

func (c *Camera) CreateBindGroup(device *wgpu.Device) error {
	var err error
	c.bindGroup, err = device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  label.Label(c, "Bind Group"),
		Layout: BindGroupLayout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding: 0,
				Buffer:  c.buffer,
				Size:    c.buffer.GetSize(),
			},
		},
	})
	return err
}

func CleanupBindGroupLayout() {
	BindGroupLayout.Release()
	BindGroupLayout = nil
}

func (c *Camera) SetBindGroup(groupIndex uint32, renderPass *wgpu.RenderPassEncoder) {
	renderPass.SetBindGroup(groupIndex, c.bindGroup, nil)
}
