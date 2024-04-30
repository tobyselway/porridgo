package instance

import (
	"unsafe"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

var VertexBufferLayout = wgpu.VertexBufferLayout{
	ArrayStride: uint64(unsafe.Sizeof(Raw{})),
	StepMode:    wgpu.VertexStepMode_Instance,
	Attributes: []wgpu.VertexAttribute{
		{
			Offset:         0,
			ShaderLocation: 5,
			Format:         wgpu.VertexFormat_Float32x4,
		},
		{
			Offset:         uint64(unsafe.Sizeof([4]float32{})),
			ShaderLocation: 6,
			Format:         wgpu.VertexFormat_Float32x4,
		},
		{
			Offset:         uint64(unsafe.Sizeof([8]float32{})),
			ShaderLocation: 7,
			Format:         wgpu.VertexFormat_Float32x4,
		},
		{
			Offset:         uint64(unsafe.Sizeof([12]float32{})),
			ShaderLocation: 8,
			Format:         wgpu.VertexFormat_Float32x4,
		},
	},
}
