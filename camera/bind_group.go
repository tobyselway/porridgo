package camera

import "github.com/rajveermalviya/go-webgpu/wgpu"

func CreateBindGroupLayout(device *wgpu.Device, label string) (*wgpu.BindGroupLayout, error) {
	return device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Label: label,
		Entries: []wgpu.BindGroupLayoutEntry{
			{
				Binding:    0,
				Visibility: wgpu.ShaderStage_Vertex,
				Buffer: wgpu.BufferBindingLayout{
					Type:             wgpu.BufferBindingType_Uniform,
					HasDynamicOffset: false,
					MinBindingSize:   0,
				},
			},
		},
	})
}

func (c *Camera) CreateBindGroup(device *wgpu.Device, layout *wgpu.BindGroupLayout, buffer *wgpu.Buffer, label string) (*wgpu.BindGroup, error) {
	return device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  label,
		Layout: layout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding: 0,
				Buffer:  buffer,
				Size:    buffer.GetSize(),
			},
		},
	})
}
