package texture

import (
	"porridgo/datatypes"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Texture struct {
	image     *datatypes.Image
	Texture   *wgpu.Texture
	View      *wgpu.TextureView
	Sampler   *wgpu.Sampler
	bindGroup *wgpu.BindGroup
}

func (t *Texture) Cleanup() {
	defer t.Texture.Release()
	defer t.View.Release()
	defer t.Sampler.Release()
	defer t.bindGroup.Release()
}
