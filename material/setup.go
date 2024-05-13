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
	if m.NormalTexture != nil {
		err := m.NormalTexture.Setup(device, queue)
		if err != nil {
			return fmt.Errorf("setting up normal texture on material %s: %w", m.Name, err)
		}
	}
	err := m.CreateBindGroup(device)
	if err != nil {
		return fmt.Errorf("creating bind group on material %s: %w", m.Name, err)
	}
	return nil
}
