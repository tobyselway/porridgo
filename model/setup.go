package model

import (
	"fmt"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (m *Model) Setup(device *wgpu.Device, queue *wgpu.Queue) error {
	for i := range m.Materials {
		err := m.Materials[i].Setup(device, queue)
		if err != nil {
			return fmt.Errorf("setting up materials for model %s: %w", m.Name, err)
		}
	}
	for i := range m.Meshes {
		err := m.Meshes[i].Setup(device, queue)
		if err != nil {
			return fmt.Errorf("setting up meshes for model %s: %w", m.Name, err)
		}
	}
	return nil
}
