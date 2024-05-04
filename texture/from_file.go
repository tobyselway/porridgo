package texture

import (
	"porridgo/datatypes"
	"porridgo/label"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func FromFile(filename string) (Texture, error) {
	image, err := datatypes.ImageFromFile(filename)
	return Texture{
		textureType: Image,
		image:       image,
	}, err
}

func (t *Texture) setupImage(device *wgpu.Device, queue *wgpu.Queue) error {
	textureExtent := wgpu.Extent3D{
		Width:              t.image.Width,
		Height:             t.image.Height,
		DepthOrArrayLayers: 1,
	}
	texture, err := device.CreateTexture(&wgpu.TextureDescriptor{
		Label:         label.Label(t, "Texture"),
		Size:          textureExtent,
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension_2D,
		Format:        wgpu.TextureFormat_RGBA8UnormSrgb,
		Usage:         wgpu.TextureUsage_TextureBinding | wgpu.TextureUsage_CopyDst,
	})
	if err != nil {
		return err
	}

	err = queue.WriteTexture(
		texture.AsImageCopy(),
		wgpu.ToBytes(t.image.Pixels[:]),
		&wgpu.TextureDataLayout{
			Offset:       0,
			BytesPerRow:  t.image.Width * 4,
			RowsPerImage: t.image.Height,
		},
		&textureExtent,
	)
	if err != nil {
		return err
	}

	view, err := texture.CreateView(nil)
	if err != nil {
		return err
	}

	sampler, err := device.CreateSampler(&wgpu.SamplerDescriptor{
		Label:          label.Label(t, "Diffuse Sampler"),
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
		return err
	}

	t.Texture = texture
	t.View = view
	t.Sampler = sampler

	return nil
}
