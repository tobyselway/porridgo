package renderer

import (
	_ "embed"
	"porridgo/camera"
	"porridgo/datatypes"
	"porridgo/texture"
	"porridgo/window"
	"unsafe"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Renderer struct {
	config            Config
	window            window.Window
	instance          *wgpu.Instance
	surface           *wgpu.Surface
	swapChain         *wgpu.SwapChain
	device            *wgpu.Device
	queue             *wgpu.Queue
	swapChainConfig   *wgpu.SwapChainDescriptor
	renderPipeline    *wgpu.RenderPipeline
	vertexBuf         *wgpu.Buffer
	indexBuf          *wgpu.Buffer
	diffuseBindGroup1 *wgpu.BindGroup
	diffuseBindGroup2 *wgpu.BindGroup
	camera            *camera.Camera
	cameraUniform     camera.Uniform
	cameraBuf         *wgpu.Buffer
	cameraBindGroup   *wgpu.BindGroup
}

func (r *Renderer) Cleanup() {
	if r.cameraBindGroup != nil {
		r.cameraBindGroup.Release()
		r.cameraBindGroup = nil
	}
	if r.cameraBuf != nil {
		r.cameraBuf.Release()
		r.cameraBuf = nil
	}
	if r.diffuseBindGroup1 != nil {
		r.diffuseBindGroup1.Release()
		r.diffuseBindGroup1 = nil
	}
	if r.diffuseBindGroup2 != nil {
		r.diffuseBindGroup2.Release()
		r.diffuseBindGroup2 = nil
	}
	if r.renderPipeline != nil {
		r.renderPipeline.Release()
		r.renderPipeline = nil
	}
	if r.indexBuf != nil {
		r.indexBuf.Release()
		r.indexBuf = nil
	}
	if r.vertexBuf != nil {
		r.vertexBuf.Release()
		r.vertexBuf = nil
	}
	if r.swapChain != nil {
		r.swapChain.Release()
		r.swapChain = nil
	}
	if r.swapChainConfig != nil {
		r.swapChainConfig = nil
	}
	if r.queue != nil {
		r.queue.Release()
		r.queue = nil
	}
	if r.device != nil {
		r.device.Release()
		r.device = nil
	}
	if r.surface != nil {
		r.surface.Release()
		r.surface = nil
	}
}

var VertexBufferLayout = wgpu.VertexBufferLayout{
	ArrayStride: uint64(unsafe.Sizeof(datatypes.Vertex{})),
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

var vertexData = [...]datatypes.Vertex{
	datatypes.NewVertex(-0.0868241, 0.49240386, 0.0, 0.4131759, 0.00759614),     // A
	datatypes.NewVertex(-0.49513406, 0.06958647, 0.0, 0.0048659444, 0.43041354), // B
	datatypes.NewVertex(-0.21918549, -0.44939706, 0.0, 0.28081453, 0.949397),    // C
	datatypes.NewVertex(0.35966998, -0.3473291, 0.0, 0.85967, 0.84732914),       // D
	datatypes.NewVertex(0.44147372, 0.2347359, 0.0, 0.9414737, 0.2652641),       // E
}

var indexData = [...]uint16{
	0, 1, 4,
	1, 2, 4,
	2, 3, 4,
}

func textureBindGroupFromPNG(device *wgpu.Device, queue *wgpu.Queue, layout *wgpu.BindGroupLayout, filename string) (*wgpu.BindGroup, error) {
	image, err := datatypes.ImageFromPNG(filename)
	if err != nil {
		return nil, err
	}

	texture, err := texture.FromImage(device, queue, image, "Diffuse Texture ("+filename+")")
	if err != nil {
		return nil, err
	}
	defer texture.Cleanup()

	bindGroup, err := texture.CreateBindGroup(device, layout, "Diffuse Bind Group ("+filename+")")
	if err != nil {
		return nil, err
	}

	return bindGroup, nil
}

//go:embed shaders/shader.wgsl
var shader string

func CreateRenderer(w window.Window, cam *camera.Camera, config Config) (r *Renderer, err error) {
	defer func() {
		if err != nil {
			r.Cleanup()
			r = nil
		}
	}()
	r = &Renderer{
		window: w,
		config: fillDefault(config),
	}

	r.camera = cam

	r.instance = wgpu.CreateInstance(nil)

	r.surface = r.window.CreateSurface(r.instance)

	adapter, err := r.instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		ForceFallbackAdapter: r.config.ForceFallbackAdapter.Unwrap(),
		CompatibleSurface:    r.surface,
	})
	if err != nil {
		return r, err
	}
	defer adapter.Release()

	r.device, err = adapter.RequestDevice(nil)
	if err != nil {
		return r, err
	}
	r.queue = r.device.GetQueue()

	caps := r.surface.GetCapabilities(adapter)

	width, height := r.window.Size()
	r.swapChainConfig = &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      caps.Formats[0],
		Width:       uint32(width),
		Height:      uint32(height),
		PresentMode: wgpu.PresentMode_Fifo,
		AlphaMode:   caps.AlphaModes[0],
	}

	r.swapChain, err = r.device.CreateSwapChain(r.surface, r.swapChainConfig)
	if err != nil {
		return r, err
	}

	r.vertexBuf, err = r.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(vertexData[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		return r, err
	}

	r.indexBuf, err = r.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Index Buffer",
		Contents: wgpu.ToBytes(indexData[:]),
		Usage:    wgpu.BufferUsage_Index,
	})
	if err != nil {
		return r, err
	}

	textureBindGroupLayout, err := texture.CreateBindGroupLayout(r.device, "Texture Bind Group Layout")
	if err != nil {
		return r, err
	}
	defer textureBindGroupLayout.Release()

	r.diffuseBindGroup1, err = textureBindGroupFromPNG(r.device, r.queue, textureBindGroupLayout, "assets/go.png")
	r.diffuseBindGroup2, err = textureBindGroupFromPNG(r.device, r.queue, textureBindGroupLayout, "assets/go2.png")

	cameraBindGroupLayout, err := camera.CreateBindGroupLayout(r.device, "Camera Bind Group Layout")
	if err != nil {
		return r, err
	}
	defer cameraBindGroupLayout.Release()

	r.cameraUniform = camera.NewUniform()
	err = r.cameraUniform.Update(r.camera)
	if err != nil {
		return r, err
	}

	r.cameraBuf, err = r.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Camera Buffer",
		Contents: wgpu.ToBytes([]camera.Uniform{r.cameraUniform}),
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		return r, err
	}

	r.cameraBindGroup, err = r.camera.CreateBindGroup(r.device, cameraBindGroupLayout, r.cameraBuf, "Camera Bind Group")
	if err != nil {
		return r, err
	}

	shader, err := r.device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label:          "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: shader},
	})
	if err != nil {
		return r, err
	}
	defer shader.Release()

	renderPipelineLayout, err := r.device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label: "Render Pipeline Layout",
		BindGroupLayouts: []*wgpu.BindGroupLayout{
			textureBindGroupLayout,
			cameraBindGroupLayout,
		},
		PushConstantRanges: []wgpu.PushConstantRange{},
	})
	if err != nil {
		return r, err
	}
	defer renderPipelineLayout.Release()

	r.renderPipeline, err = r.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label:  "Render Pipeline",
		Layout: renderPipelineLayout,
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
					Format:    r.swapChainConfig.Format,
					Blend:     &wgpu.BlendState_Replace,
					WriteMask: wgpu.ColorWriteMask_All,
				},
			},
		},
	})
	if err != nil {
		return r, err
	}

	r.window.OnResize(func(width int, height int) {
		r.Resize(width, height)
	})

	return r, nil
}

func (r *Renderer) Resize(width int, height int) {
	if width > 0 && height > 0 {
		r.swapChainConfig.Width = uint32(width)
		r.swapChainConfig.Height = uint32(height)

		if r.swapChain != nil {
			r.swapChain.Release()
		}
		var err error
		r.swapChain, err = r.device.CreateSwapChain(r.surface, r.swapChainConfig)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Renderer) Render(spacePressed bool) error {
	nextTexture, err := r.swapChain.GetCurrentTextureView()
	if err != nil {
		return err
	}
	defer nextTexture.Release()

	encoder, err := r.device.CreateCommandEncoder(nil)
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

	renderPass.SetPipeline(r.renderPipeline)

	if spacePressed {
		renderPass.SetBindGroup(0, r.diffuseBindGroup2, nil)
	} else {
		renderPass.SetBindGroup(0, r.diffuseBindGroup1, nil)
	}

	r.cameraUniform.Update(r.camera)
	r.queue.WriteBuffer(r.cameraBuf, 0, wgpu.ToBytes([]camera.Uniform{r.cameraUniform}))

	renderPass.SetBindGroup(1, r.cameraBindGroup, nil)

	renderPass.SetVertexBuffer(0, r.vertexBuf, 0, wgpu.WholeSize)
	renderPass.SetIndexBuffer(r.indexBuf, wgpu.IndexFormat_Uint16, 0, wgpu.WholeSize)

	renderPass.DrawIndexed(uint32(len(indexData)), 1, 0, 0, 0)
	renderPass.End()

	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		return err
	}
	defer cmdBuffer.Release()

	r.queue.Submit(cmdBuffer)
	r.swapChain.Present()

	return nil
}
