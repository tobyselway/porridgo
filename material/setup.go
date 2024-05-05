package material

import (
	"fmt"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (m *Material) Setup(device *wgpu.Device, queue *wgpu.Queue) error {
	if m.DiffuseTexture != nil {
		err := m.DiffuseTexture.Setup(device, queue)
		if err != nil {
			return fmt.Errorf("setting up diffuse texture on material %s: %w", m.Name, err)
		}
	}
	return nil
}
