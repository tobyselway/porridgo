package renderer

import (
	"embed"
	_ "embed"
	"porridgo/camera"
	"porridgo/datatypes"
	"porridgo/instance"
	"porridgo/light"
	"porridgo/mesh"
	"porridgo/model"
	"porridgo/pipeline"
	"porridgo/shader"
	"porridgo/texture"
	"porridgo/window"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Renderer struct {
	config              Config
	window              window.Window
	wgpuInstance        *wgpu.Instance
	surface             *wgpu.Surface
	swapChain           *wgpu.SwapChain
	device              *wgpu.Device
	queue               *wgpu.Queue
	swapChainConfig     *wgpu.SwapChainDescriptor
	renderPipeline      *pipeline.Pipeline
	lightRenderPipeline *pipeline.Pipeline
	depthTexture        texture.Texture
	camera              *camera.Camera
	instances           []instance.Instance
	instanceBuf         *wgpu.Buffer
	mdl                 *model.Model
	sun                 *light.Light
}

func (r *Renderer) Cleanup() {
	if r.camera != nil {
		r.camera.Cleanup()
		r.camera = nil
	}
	if r.sun != nil {
		r.sun.Cleanup()
		r.sun = nil
	}
	if r.mdl != nil {
		r.mdl.Cleanup()
		r.mdl = nil
	}
	if r.renderPipeline != nil {
		r.renderPipeline.Cleanup()
		r.renderPipeline = nil
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

// Any fs.FS filesystem can be provided to load shader files from.
// In this case an embed.FS is used to load shaders at compile-time.
//
//go:embed shaders
var embedShaders embed.FS

const NUM_INSTANCES_PER_ROW uint32 = 10

var SPACE_BETWEEN float32 = 3.0

func GenerateInstances() []instance.Instance {
	instances := []instance.Instance{}

	for i := uint32(0); i < NUM_INSTANCES_PER_ROW; i++ {
		for j := uint32(0); j < NUM_INSTANCES_PER_ROW; j++ {
			x := SPACE_BETWEEN * (float32(i) - float32(NUM_INSTANCES_PER_ROW)/2.0)
			z := SPACE_BETWEEN * (float32(j) - float32(NUM_INSTANCES_PER_ROW)/2.0)
			position := datatypes.NewVec3f(x, 0.0, z)
			instances = append(instances, instance.Instance{
				Position: position,
				Scale:    datatypes.NewVec3f(1.0, 1.0, 1.0),
				Rotation: datatypes.NewVec3f(0.0, 0.0, 0.0),
			})
		}
	}
	return instances
}

func CreateRenderer(w window.Window, cam *camera.Camera, sun *light.Light, mdl *model.Model, config Config) (r *Renderer, err error) {
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
	r.sun = sun
	r.mdl = mdl

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

	r.instances = GenerateInstances()

	instanceData := []instance.Raw{
		// instance.Instance{
		// 	Position: datatypes.NewVec3f(0.0, 0.0, 0.0),
		// 	Scale:    datatypes.NewVec3f(1.0, 1.0, 1.0),
		// 	Rotation: datatypes.NewVec3f(0.0, 0.0, 0.0),
		// }.ToRaw(),
	}
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

	err = light.SetupBindGroupLayout(r.device)
	if err != nil {
		return r, err
	}
	defer light.CleanupBindGroupLayout()

	err = r.mdl.Setup(r.device, r.queue)
	if err != nil {
		return r, err
	}

	err = r.camera.Setup(r.device, r.queue)
	if err != nil {
		return r, err
	}

	err = r.sun.Setup(r.device, r.queue)
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

	shd, err := shader.FromFile("shaders/shader.wgsl", embedShaders)
	if err != nil {
		return r, err
	}

	r.renderPipeline = &pipeline.Pipeline{
		Name:           "Render",
		VertexShader:   shd,
		FragmentShader: shd,
		BindGroupLayouts: []*wgpu.BindGroupLayout{
			texture.BindGroupLayout,
			camera.BindGroupLayout,
			light.BindGroupLayout,
		},
		VertexBufferLayouts: []wgpu.VertexBufferLayout{
			mesh.VertexBufferLayout,
			instance.VertexBufferLayout,
		},
	}

	err = r.renderPipeline.Setup(r.device, r.swapChainConfig)
	if err != nil {
		return r, err
	}

	shdLight, err := shader.FromFile("shaders/light.wgsl", embedShaders)
	if err != nil {
		return r, err
	}

	r.lightRenderPipeline = &pipeline.Pipeline{
		Name:           "Light",
		VertexShader:   shdLight,
		FragmentShader: shdLight,
		BindGroupLayouts: []*wgpu.BindGroupLayout{
			camera.BindGroupLayout,
			light.BindGroupLayout,
		},
		VertexBufferLayouts: []wgpu.VertexBufferLayout{
			mesh.VertexBufferLayout,
		},
	}

	err = r.lightRenderPipeline.Setup(r.device, r.swapChainConfig)
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
				ClearValue: wgpu.Color{R: 0.1, G: 0.1, B: 0.1, A: 1.0},
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

	r.renderPipeline.Prepare(renderPass)
	r.camera.Prepare(1, r.queue, renderPass)
	r.sun.Prepare(2, r.queue, renderPass)
	renderPass.SetVertexBuffer(1, r.instanceBuf, 0, wgpu.WholeSize)
	r.mdl.DrawInstanced(uint32(len(r.instances)), renderPass)

	r.lightRenderPipeline.Prepare(renderPass)
	r.camera.Prepare(0, r.queue, renderPass)
	r.sun.Prepare(1, r.queue, renderPass)
	r.sun.Draw(r.mdl, renderPass)

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
