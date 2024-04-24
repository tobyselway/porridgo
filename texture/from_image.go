package texture

import (
	"porridgo/datatypes"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func FromImage(device *wgpu.Device, queue *wgpu.Queue, image *datatypes.Image, label string) (*Texture, error) {
	textureExtent := wgpu.Extent3D{
		Width:              image.Width,
		Height:             image.Height,
		DepthOrArrayLayers: 1,
	}
	texture, err := device.CreateTexture(&wgpu.TextureDescriptor{
		Label:         label,
		Size:          textureExtent,
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension_2D,
		Format:        wgpu.TextureFormat_RGBA8UnormSrgb,
		Usage:         wgpu.TextureUsage_TextureBinding | wgpu.TextureUsage_CopyDst,
	})
	if err != nil {
		return nil, err
	}

	err = queue.WriteTexture(
		texture.AsImageCopy(),
		wgpu.ToBytes(image.Pixels[:]),
		&wgpu.TextureDataLayout{
			Offset:       0,
			BytesPerRow:  image.Width * 4,
			RowsPerImage: image.Height,
		},
		&textureExtent,
	)
	if err != nil {
		return nil, err
	}

	view, err := texture.CreateView(nil)
	if err != nil {
		return nil, err
	}

	sampler, err := device.CreateSampler(&wgpu.SamplerDescriptor{
		Label:          "Diffuse Sampler",
		AddressModeU:   wgpu.AddressMode_ClampToEdge,
		AddressModeV:   wgpu.AddressMode_ClampToEdge,
		AddressModeW:   wgpu.AddressMode_ClampToEdge,
		MagFilter:      wgpu.FilterMode_Linear,
		MinFilter:      wgpu.FilterMode_Nearest,
		MipmapFilter:   wgpu.MipmapFilterMode_Nearest,
		LodMinClamp:    0.0,
		LodMaxClamp:    32.0,
		Compare:        wgpu.CompareFunction_Undefined,
		MaxAnisotrophy: 1,
	})
	if err != nil {
		return nil, err
	}

	return &Texture{
		Texture: texture,
		View:    view,
		Sampler: sampler,
	}, nil
}
