package material

import (
	"porridgo/texture"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Material struct {
	Name           string
	DiffuseTexture *texture.Texture
	NormalTexture  *texture.Texture
	bindGroup      *wgpu.BindGroup
}

func (m *Material) Cleanup() {
	if m.DiffuseTexture != nil {
		defer m.DiffuseTexture.Cleanup()
	}
	if m.NormalTexture != nil {
		defer m.NormalTexture.Cleanup()
	}
	if m.bindGroup != nil {
		defer m.bindGroup.Release()
	}
}
