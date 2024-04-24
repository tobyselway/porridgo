package texture

import (
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Texture struct {
	Texture *wgpu.Texture
	View    *wgpu.TextureView
	Sampler *wgpu.Sampler
}

func (t *Texture) Cleanup() {
	defer t.Texture.Release()
	defer t.View.Release()
	defer t.Sampler.Release()
}
