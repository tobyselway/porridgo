package main

import (
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
	wgpuext_glfw "github.com/rajveermalviya/go-webgpu/wgpuext/glfw"

	_ "embed"
)

var forceFallbackAdapter = os.Getenv("WGPU_FORCE_FALLBACK_ADAPTER") == "1"

func init() {
	runtime.LockOSThread()

	switch os.Getenv("WGPU_LOG_LEVEL") {
	case "OFF":
		wgpu.SetLogLevel(wgpu.LogLevel_Off)
	case "ERROR":
		wgpu.SetLogLevel(wgpu.LogLevel_Error)
	case "WARN":
		wgpu.SetLogLevel(wgpu.LogLevel_Warn)
	case "INFO":
		wgpu.SetLogLevel(wgpu.LogLevel_Info)
	case "DEBUG":
		wgpu.SetLogLevel(wgpu.LogLevel_Debug)
	case "TRACE":
		wgpu.SetLogLevel(wgpu.LogLevel_Trace)
	}
}

//go:embed assets/shader.wgsl
var shader string

type State struct {
	instance  *wgpu.Instance
	surface   *wgpu.Surface
	swapChain *wgpu.SwapChain
	device    *wgpu.Device
	queue     *wgpu.Queue
	config    *wgpu.SwapChainDescriptor
	pipeline  *wgpu.RenderPipeline
	vertexBuf *wgpu.Buffer
	indexBuf  *wgpu.Buffer
	bindGroup *wgpu.BindGroup
}

type Vertex struct {
	pos      [3]float32
	texCoord [2]float32
}

var VertexBufferLayout = wgpu.VertexBufferLayout{
	ArrayStride: uint64(unsafe.Sizeof(Vertex{})),
	StepMode:    wgpu.VertexStepMode_Vertex,
	Attributes: []wgpu.VertexAttribute{
		{
			Format:         wgpu.VertexFormat_Float32x3,
			Offset:         0,
			ShaderLocation: 0,
		},
		{
			Format:         wgpu.VertexFormat_Float32x2,
			Offset:         uint64(unsafe.Sizeof([3]float32{})),
			ShaderLocation: 1,
		},
	},
}

func vertex(pos1, pos2, pos3, tc1, tc2 float32) Vertex {
	return Vertex{
		pos:      [3]float32{pos1, pos2, pos3},
		texCoord: [2]float32{tc1, tc2},
	}
}

var vertexData = [...]Vertex{
	vertex(-1, -1, 0, 1, 1), // TL
	vertex(1, -1, 0, 0, 1),  // TR
	vertex(-1, 1, 0, 1, 0),  // BL
	vertex(1, 1, 0, 0, 0),   // BR
}

var indexData = [...]uint16{
	0, 2, 1,
	1, 2, 3,
}

const texelsSize = 256

func createTexels() (texels [(texelsSize * 4) * texelsSize]uint8) {

	imagePath := "assets/go.png"
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		return
	}
	defer file.Close()

	// Decode the PNG image
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("Error decoding PNG:", err)
		return
	}

	// Populate pixel array
	for y := 0; y < texelsSize; y++ {
		for x := 0; x < texelsSize; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			texels[y*4*texelsSize+x*4+0] = uint8(r >> 8)
			texels[y*4*texelsSize+x*4+1] = uint8(g >> 8)
			texels[y*4*texelsSize+x*4+2] = uint8(b >> 8)
			texels[y*4*texelsSize+x*4+3] = 0xff
		}
	}

	return texels
}

func InitState(window *glfw.Window) (s *State, err error) {
	defer func() {
		if err != nil {
			s.Destroy()
			s = nil
		}
	}()
	s = &State{}

	s.instance = wgpu.CreateInstance(nil)

	s.surface = s.instance.CreateSurface(wgpuext_glfw.GetSurfaceDescriptor(window))

	adapter, err := s.instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		ForceFallbackAdapter: forceFallbackAdapter,
		CompatibleSurface:    s.surface,
	})
	if err != nil {
		return s, err
	}
	defer adapter.Release()

	s.device, err = adapter.RequestDevice(nil)
	if err != nil {
		return s, err
	}
	s.queue = s.device.GetQueue()

	caps := s.surface.GetCapabilities(adapter)

	width, height := window.GetSize()
	s.config = &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      caps.Formats[0],
		Width:       uint32(width),
		Height:      uint32(height),
		PresentMode: wgpu.PresentMode_Fifo,
		AlphaMode:   caps.AlphaModes[0],
	}

	s.swapChain, err = s.device.CreateSwapChain(s.surface, s.config)
	if err != nil {
		return s, err
	}

	s.vertexBuf, err = s.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(vertexData[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		return s, err
	}

	s.indexBuf, err = s.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Index Buffer",
		Contents: wgpu.ToBytes(indexData[:]),
		Usage:    wgpu.BufferUsage_Index,
	})
	if err != nil {
		return s, err
	}

	texels := createTexels()
	textureExtent := wgpu.Extent3D{
		Width:              texelsSize,
		Height:             texelsSize,
		DepthOrArrayLayers: 1,
	}
	texture, err := s.device.CreateTexture(&wgpu.TextureDescriptor{
		Size:          textureExtent,
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension_2D,
		Format:        wgpu.TextureFormat_RGBA8UnormSrgb,
		Usage:         wgpu.TextureUsage_TextureBinding | wgpu.TextureUsage_CopyDst,
	})
	if err != nil {
		return s, err
	}
	defer texture.Release()

	textureView, err := texture.CreateView(nil)
	if err != nil {
		return s, err
	}
	defer textureView.Release()

	s.queue.WriteTexture(
		texture.AsImageCopy(),
		wgpu.ToBytes(texels[:]),
		&wgpu.TextureDataLayout{
			Offset:       0,
			BytesPerRow:  texelsSize * 4,
			RowsPerImage: wgpu.CopyStrideUndefined,
		},
		&textureExtent,
	)

	if err != nil {
		return s, err
	}

	shader, err := s.device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label:          "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: shader},
	})
	if err != nil {
		return s, err
	}
	defer shader.Release()

	s.pipeline, err = s.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label: "Render Pipeline",
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
			Buffers: []wgpu.VertexBufferLayout{
				VertexBufferLayout,
			},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:         wgpu.PrimitiveTopology_TriangleList,
			StripIndexFormat: wgpu.IndexFormat_Undefined,
			FrontFace:        wgpu.FrontFace_CCW,
			CullMode:         wgpu.CullMode_None,
		},
		Multisample: wgpu.MultisampleState{
			Count:                  1,
			Mask:                   0xFFFFFFFF,
			AlphaToCoverageEnabled: false,
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{
				{
					Format:    s.config.Format,
					Blend:     &wgpu.BlendState_Replace,
					WriteMask: wgpu.ColorWriteMask_All,
				},
			},
		},
	})
	if err != nil {
		return s, err
	}

	bindGroupLayout := s.pipeline.GetBindGroupLayout(0)
	defer bindGroupLayout.Release()

	s.bindGroup, err = s.device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Layout: bindGroupLayout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding:     0,
				TextureView: textureView,
				Size:        wgpu.WholeSize,
			},
		},
	})
	if err != nil {
		return s, err
	}

	return s, nil
}

func (s *State) Resize(width, height int) {
	if width > 0 && height > 0 {
		s.config.Width = uint32(width)
		s.config.Height = uint32(height)

		if s.swapChain != nil {
			s.swapChain.Release()
		}
		var err error
		s.swapChain, err = s.device.CreateSwapChain(s.surface, s.config)
		if err != nil {
			panic(err)
		}
	}
}

func (s *State) Render() error {
	nextTexture, err := s.swapChain.GetCurrentTextureView()
	if err != nil {
		return err
	}
	defer nextTexture.Release()

	encoder, err := s.device.CreateCommandEncoder(nil)
	if err != nil {
		return err
	}
	defer encoder.Release()

	renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{
			{
				View:       nextTexture,
				LoadOp:     wgpu.LoadOp_Clear,
				StoreOp:    wgpu.StoreOp_Store,
				ClearValue: wgpu.Color{R: 0.2, G: 0.2, B: 0.2, A: 1.0},
			},
		},
	})
	defer renderPass.Release()

	renderPass.SetPipeline(s.pipeline)
	renderPass.SetBindGroup(0, s.bindGroup, nil)
	renderPass.SetIndexBuffer(s.indexBuf, wgpu.IndexFormat_Uint16, 0, wgpu.WholeSize)
	renderPass.SetVertexBuffer(0, s.vertexBuf, 0, wgpu.WholeSize)
	renderPass.DrawIndexed(uint32(len(indexData)), 1, 0, 0, 0)
	renderPass.End()

	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		return err
	}
	defer cmdBuffer.Release()

	s.queue.Submit(cmdBuffer)
	s.swapChain.Present()

	return nil
}

func (s *State) Destroy() {
	if s.bindGroup != nil {
		s.bindGroup.Release()
		s.bindGroup = nil
	}
	if s.pipeline != nil {
		s.pipeline.Release()
		s.pipeline = nil
	}
	if s.indexBuf != nil {
		s.indexBuf.Release()
		s.indexBuf = nil
	}
	if s.vertexBuf != nil {
		s.vertexBuf.Release()
		s.vertexBuf = nil
	}
	if s.swapChain != nil {
		s.swapChain.Release()
		s.swapChain = nil
	}
	if s.config != nil {
		s.config = nil
	}
	if s.queue != nil {
		s.queue.Release()
		s.queue = nil
	}
	if s.device != nil {
		s.device.Release()
		s.device = nil
	}
	if s.surface != nil {
		s.surface.Release()
		s.surface = nil
	}
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	window, err := glfw.CreateWindow(640, 480, "porridgo", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	s, err := InitState(window)
	if err != nil {
		panic(err)
	}
	defer s.Destroy()

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// Print resource usage on pressing 'R'
		if key == glfw.KeyR && (action == glfw.Press || action == glfw.Repeat) {
			report := s.instance.GenerateReport()
			buf, _ := json.MarshalIndent(report, "", "  ")
			fmt.Print(string(buf))
		}
	})

	window.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		s.Resize(width, height)
	})

	for !window.ShouldClose() {
		glfw.PollEvents()

		err := s.Render()
		if err != nil {
			fmt.Println("error occured while rendering:", err)

			errstr := err.Error()
			switch {
			case strings.Contains(errstr, "Surface timed out"): // do nothing
			case strings.Contains(errstr, "Surface is outdated"): // do nothing
			case strings.Contains(errstr, "Surface was lost"): // do nothing
			default:
				panic(err)
			}
		}
	}
}
