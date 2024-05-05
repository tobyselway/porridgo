package material

import (
	"porridgo/texture"
)

type Material struct {
	Name           string
	DiffuseTexture *texture.Texture
}

func (m *Material) Cleanup() {
	if m.DiffuseTexture != nil {
		defer m.DiffuseTexture.Cleanup()
	}
}
