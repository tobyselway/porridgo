package texture

import "github.com/rajveermalviya/go-webgpu/wgpu"

func CreateBindGroupLayout(device *wgpu.Device, label string) (*wgpu.BindGroupLayout, error) {
	return device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Label: label,
		Entries: []wgpu.BindGroupLayoutEntry{
			{
				Binding:    0,
				Visibility: wgpu.ShaderStage_Fragment,
				Texture: wgpu.TextureBindingLayout{
					Multisampled:  false,
					ViewDimension: wgpu.TextureViewDimension_2D,
					SampleType:    wgpu.TextureSampleType_Float,
				},
			},
			{
				Binding:    1,
				Visibility: wgpu.ShaderStage_Fragment,
				Sampler: wgpu.SamplerBindingLayout{
					Type: wgpu.SamplerBindingType_Filtering,
				},
			},
		},
	})
}

func (t *Texture) CreateBindGroup(device *wgpu.Device, layout *wgpu.BindGroupLayout, label string) (*wgpu.BindGroup, error) {
	return device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  label,
		Layout: layout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding:     0,
				TextureView: t.View,
			},
			{
				Binding: 1,
				Sampler: t.Sampler,
			},
		},
	})
}
