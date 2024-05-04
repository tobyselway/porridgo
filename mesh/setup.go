package mesh

import (
	"fmt"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (m *Mesh) Setup(device *wgpu.Device, queue *wgpu.Queue) error {
	err := m.SetupBuffers(device)
	if err != nil {
		return fmt.Errorf("setting up buffers on mesh %s: %w", m.Name, err)
	}
	return nil
}
