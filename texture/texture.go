package texture

import (
	"porridgo/datatypes"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Type int

const (
	Image Type = iota
	Depth
)

type Texture struct {
	textureType Type
	image       *datatypes.Image
	Texture     *wgpu.Texture
	View        *wgpu.TextureView
	Sampler     *wgpu.Sampler
	depthConfig *DepthConfig
}

func (t *Texture) Cleanup() {
	if t.Texture != nil {
		defer t.Texture.Release()
	}
	if t.View != nil {
		defer t.View.Release()
	}
	if t.Sampler != nil {
		defer t.Sampler.Release()
	}
}
