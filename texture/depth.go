package texture

import (
	"porridgo/label"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

const DEPTH_FORMAT = wgpu.TextureFormat_Depth32Float

type DepthConfig struct {
	Width  uint32
	Height uint32
}

func CreateDepthTexture(config *DepthConfig) Texture {
	return Texture{
		textureType: Depth,
		depthConfig: config,
	}
}

func (t *Texture) setupDepth(device *wgpu.Device) error {
	textureExtent := wgpu.Extent3D{
		Width:              t.depthConfig.Width,
		Height:             t.depthConfig.Height,
		DepthOrArrayLayers: 1,
	}
	texture, err := device.CreateTexture(&wgpu.TextureDescriptor{
		Label:         label.Label(t, "Texture"),
		Size:          textureExtent,
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension_2D,
		Format:        DEPTH_FORMAT,
		Usage:         wgpu.TextureUsage_RenderAttachment | wgpu.TextureUsage_TextureBinding,
	})
	if err != nil {
		return err
	}

	view, err := texture.CreateView(nil)
	if err != nil {
		return err
	}

	sampler, err := device.CreateSampler(&wgpu.SamplerDescriptor{
		Label:          label.Label(t, "Sampler"),
		AddressModeU:   wgpu.AddressMode_ClampToEdge,
		AddressModeV:   wgpu.AddressMode_ClampToEdge,
		AddressModeW:   wgpu.AddressMode_ClampToEdge,
		MagFilter:      wgpu.FilterMode_Linear,
		MinFilter:      wgpu.FilterMode_Linear,
		MipmapFilter:   wgpu.MipmapFilterMode_Nearest,
		Compare:        wgpu.CompareFunction_LessEqual,
		LodMinClamp:    0.0,
		LodMaxClamp:    100.0,
		MaxAnisotrophy: 1,
	})
	if err != nil {
		return err
	}

	t.Texture = texture
	t.View = view
	t.Sampler = sampler

	return nil
}
