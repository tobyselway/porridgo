package texture

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
		Label: "Texture Bind Group Layout",
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
	return err
}

func CleanupBindGroupLayout() {
	BindGroupLayout.Release()
	BindGroupLayout = nil
}

func (t *Texture) CreateBindGroup(device *wgpu.Device) error {
	var err error
	t.bindGroup, err = device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  label.Label(t, "Bind Group"),
		Layout: BindGroupLayout,
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
	return err
}

func (t *Texture) SetBindGroup(renderPass *wgpu.RenderPassEncoder) {
	renderPass.SetBindGroup(0, t.bindGroup, nil)
}
