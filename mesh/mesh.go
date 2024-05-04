package mesh

import "github.com/rajveermalviya/go-webgpu/wgpu"

type Mesh struct {
	Name         string
	VertexBuffer *wgpu.Buffer
	IndexBuffer  *wgpu.Buffer
	NumElements  uint32
	Material     uint32
	Vertices     []Vertex
	Indices      []uint32
}

func (m *Mesh) Cleanup() {
	if m.VertexBuffer != nil {
		defer m.VertexBuffer.Destroy()
	}
	if m.IndexBuffer != nil {
		defer m.IndexBuffer.Destroy()
	}
}
