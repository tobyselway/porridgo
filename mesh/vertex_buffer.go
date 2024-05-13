package mesh

import (
	"unsafe"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

var VertexBufferLayout = wgpu.VertexBufferLayout{
	ArrayStride: uint64(unsafe.Sizeof(Vertex{})),
	StepMode:    wgpu.VertexStepMode_Vertex,
	Attributes: []wgpu.VertexAttribute{
		// Position
		{
			Offset:         0,
			ShaderLocation: 0,
			Format:         wgpu.VertexFormat_Float32x3,
		},
		// TexCoords
		{
			Offset:         uint64(unsafe.Sizeof([3]float32{})),
			ShaderLocation: 1,
			Format:         wgpu.VertexFormat_Float32x2,
		},
		// Normal
		{
			Offset:         uint64(unsafe.Sizeof([5]float32{})),
			ShaderLocation: 2,
			Format:         wgpu.VertexFormat_Float32x3,
		},
		// Tangent
		{
			Offset:         uint64(unsafe.Sizeof([8]float32{})),
			ShaderLocation: 3,
			Format:         wgpu.VertexFormat_Float32x3,
		},
		// Bitangent
		{
			Offset:         uint64(unsafe.Sizeof([11]float32{})),
			ShaderLocation: 4,
			Format:         wgpu.VertexFormat_Float32x3,
		},
	},
}
