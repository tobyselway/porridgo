package material

import (
	"porridgo/texture"
)

type Material struct {
	Name           string
	DiffuseTexture texture.Texture
}

func (m *Material) Cleanup() {
	defer m.DiffuseTexture.Cleanup()
}
