package mesh

import (
	"fmt"
	"porridgo/label"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (m *Mesh) SetupBuffers(device *wgpu.Device) error {
	var err error
	m.VertexBuffer, err = device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    label.Label(m, "Vertex Buffer"),
		Contents: wgpu.ToBytes(m.Vertices[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		return fmt.Errorf("creating vertex buffer for %s mesh: %w", m.Name, err)
	}

	m.IndexBuffer, err = device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    label.Label(m, "Index Buffer"),
		Contents: wgpu.ToBytes(m.Indices[:]),
		Usage:    wgpu.BufferUsage_Index,
	})
	if err != nil {
		return fmt.Errorf("creating index buffer for %s mesh: %w", m.Name, err)
	}
	return err
}
