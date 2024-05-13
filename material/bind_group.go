package material

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
		Label: "Material Bind Group Layout",
		Entries: []wgpu.BindGroupLayoutEntry{
			// Diffuse
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
			// Normal
			{
				Binding:    2,
				Visibility: wgpu.ShaderStage_Fragment,
				Texture: wgpu.TextureBindingLayout{
					Multisampled:  false,
					ViewDimension: wgpu.TextureViewDimension_2D,
					SampleType:    wgpu.TextureSampleType_Float,
				},
			},
			{
				Binding:    3,
				Visibility: wgpu.ShaderStage_Fragment,
				Sampler: wgpu.SamplerBindingLayout{
					Type: wgpu.SamplerBindingType_Filtering,
				},
			},
		},
	})
	return err
}

func CleanupBindGroupLayout() {
	BindGroupLayout.Release()
	BindGroupLayout = nil
}

func (m *Material) CreateBindGroup(device *wgpu.Device) error {
	var err error
	m.bindGroup, err = device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  label.Label(m, "Bind Group"),
		Layout: BindGroupLayout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding:     0,
				TextureView: m.DiffuseTexture.View,
			},
			{
				Binding: 1,
				Sampler: m.DiffuseTexture.Sampler,
			},
			{
				Binding:     2,
				TextureView: m.NormalTexture.View,
			},
			{
				Binding: 3,
				Sampler: m.NormalTexture.Sampler,
			},
		},
	})
	return err
}

func (m *Material) SetBindGroup(renderPass *wgpu.RenderPassEncoder) {
	renderPass.SetBindGroup(0, m.bindGroup, nil)
}
