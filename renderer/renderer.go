package renderer

import (
	_ "embed"
	"porridgo/camera"
	"porridgo/datatypes"
	"porridgo/instance"
	"porridgo/texture"
	"porridgo/vertex"
	"porridgo/window"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Renderer struct {
	config          Config
	window          window.Window
	wgpuInstance    *wgpu.Instance
	surface         *wgpu.Surface
	swapChain       *wgpu.SwapChain
	device          *wgpu.Device
	queue           *wgpu.Queue
	swapChainConfig *wgpu.SwapChainDescriptor
	renderPipeline  *wgpu.RenderPipeline
	vertexBuf       *wgpu.Buffer
	indexBuf        *wgpu.Buffer
	texture1        *texture.Texture
	texture2        *texture.Texture
	depthTexture    texture.Texture
	camera          *camera.Camera
	instances       []instance.Instance
	instanceBuf     *wgpu.Buffer
}

func (r *Renderer) Cleanup() {
	if r.camera != nil {
		r.camera.Cleanup()
		r.camera = nil
	}
	if r.texture1 != nil {
		r.texture1.Cleanup()
		r.texture1 = nil
	}
	if r.texture2 != nil {
		r.texture2.Cleanup()
		r.texture2 = nil
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

var vertexData = [...]vertex.Vertex{
	vertex.NewVertex(-0.0868241, 0.49240386, 0.0, 0.4131759, 0.00759614),     // A
	vertex.NewVertex(-0.49513406, 0.06958647, 0.0, 0.0048659444, 0.43041354), // B
	vertex.NewVertex(-0.21918549, -0.44939706, 0.0, 0.28081453, 0.949397),    // C
	vertex.NewVertex(0.35966998, -0.3473291, 0.0, 0.85967, 0.84732914),       // D
	vertex.NewVertex(0.44147372, 0.2347359, 0.0, 0.9414737, 0.2652641),       // E
}

var indexData = [...]uint16{
	0, 1, 4,
	1, 2, 4,
	2, 3, 4,
}

//go:embed shaders/shader.wgsl
var shader string

const NUM_INSTANCES_PER_ROW uint32 = 10

var INSTANCE_DISPLACEMENT datatypes.Vec3f = datatypes.NewVec3f(
	float32(NUM_INSTANCES_PER_ROW)*0.5,
	0.0,
	float32(NUM_INSTANCES_PER_ROW)*0.5,
)

func GenerateInstances() []instance.Instance {
	instances := []instance.Instance{}
	for z := 0; z < int(NUM_INSTANCES_PER_ROW); z++ {
		for x := 0; x < int(NUM_INSTANCES_PER_ROW); x++ {
			position := datatypes.NewVec3f(float32(x), 0.0, float32(z)).Sub(INSTANCE_DISPLACEMENT)
			instances = append(instances, instance.Instance{
				Position: position,
				Scale:    datatypes.NewVec3f(1.0, 1.0, 1.0),
				Rotation: datatypes.NewVec3f(0.0, 0.0, 0.0),
			})
		}
	}
	return instances
}

func CreateRenderer(w window.Window, cam *camera.Camera, tex1 *texture.Texture, tex2 *texture.Texture, config Config) (r *Renderer, err error) {
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
	r.texture1 = tex1
	r.texture2 = tex2

	r.wgpuInstance = wgpu.CreateInstance(nil)

	r.surface = r.window.CreateSurface(r.wgpuInstance)

	adapter, err := r.wgpuInstance.RequestAdapter(&wgpu.RequestAdapterOptions{
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

	r.instances = GenerateInstances()

	instanceData := []instance.Raw{}
	for _, inst := range r.instances {
		instanceData = append(instanceData, inst.ToRaw())
	}

	r.instanceBuf, err = r.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Instance Buffer",
		Contents: wgpu.ToBytes(instanceData[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		return r, err
	}

	err = texture.SetupBindGroupLayout(r.device)
	if err != nil {
		return r, err
	}
	defer texture.CleanupBindGroupLayout()

	err = camera.SetupBindGroupLayout(r.device)
	if err != nil {
		return r, err
	}
	defer camera.CleanupBindGroupLayout()

	err = r.texture1.Setup(r.device, r.queue)
	if err != nil {
		return r, err
	}
	err = r.texture1.CreateBindGroup(r.device)
	if err != nil {
		return r, err
	}

	err = r.texture2.Setup(r.device, r.queue)
	if err != nil {
		return r, err
	}
	err = r.texture2.CreateBindGroup(r.device)
	if err != nil {
		return r, err
	}

	r.camera.SetupUniform()
	err = r.camera.UpdateUniform()
	if err != nil {
		return r, err
	}

	err = r.camera.SetupBuffer(r.device)
	if err != nil {
		return r, err
	}

	err = r.camera.CreateBindGroup(r.device)
	if err != nil {
		return r, err
	}

	r.depthTexture = texture.CreateDepthTexture(&texture.DepthConfig{
		Width:  uint32(width),
		Height: uint32(height),
	})

	err = r.depthTexture.Setup(r.device, r.queue)
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
			texture.BindGroupLayout,
			camera.BindGroupLayout,
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
				vertex.VertexBufferLayout,
				instance.VertexBufferLayout,
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
		DepthStencil: &wgpu.DepthStencilState{
			Format:            texture.DEPTH_FORMAT,
			DepthWriteEnabled: true,
			DepthCompare:      wgpu.CompareFunction_Less,
			StencilFront: wgpu.StencilFaceState{
				Compare: wgpu.CompareFunction_Never,
			},
			StencilBack: wgpu.StencilFaceState{
				Compare: wgpu.CompareFunction_Never,
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

		r.depthTexture = texture.CreateDepthTexture(&texture.DepthConfig{
			Width:  uint32(width),
			Height: uint32(height),
		})
		err = r.depthTexture.Setup(r.device, r.queue)
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
		DepthStencilAttachment: &wgpu.RenderPassDepthStencilAttachment{
			View:            r.depthTexture.View,
			DepthLoadOp:     wgpu.LoadOp_Clear,
			DepthClearValue: 1.0,
			DepthStoreOp:    wgpu.StoreOp_Store,
			StencilLoadOp:   wgpu.LoadOp_Load,
			StencilStoreOp:  wgpu.StoreOp_Discard,
		},
	})
	defer renderPass.Release()

	renderPass.SetPipeline(r.renderPipeline)

	// for renderable := range renderables {
	// 	renderable.UpdateUniforms(r.queue)
	// 	renderable.SetBindGroups(renderPass)
	// }

	if spacePressed {
		r.texture2.SetBindGroup(renderPass)
	} else {
		r.texture1.SetBindGroup(renderPass)
	}

	r.camera.UpdateUniform()
	r.camera.WriteBuffer(r.queue)
	r.camera.SetBindGroup(renderPass)

	renderPass.SetVertexBuffer(0, r.vertexBuf, 0, wgpu.WholeSize)
	renderPass.SetVertexBuffer(1, r.instanceBuf, 0, wgpu.WholeSize)
	renderPass.SetIndexBuffer(r.indexBuf, wgpu.IndexFormat_Uint16, 0, wgpu.WholeSize)

	renderPass.DrawIndexed(uint32(len(indexData)), uint32(len(r.instances)), 0, 0, 0)
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
