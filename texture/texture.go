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
	bindGroup   *wgpu.BindGroup
	depthConfig *DepthConfig
}

func (t *Texture) Cleanup() {
	defer t.Texture.Release()
	defer t.View.Release()
	defer t.Sampler.Release()
	defer t.bindGroup.Release()
}
