package glfw

import (
	"github.com/rajveermalviya/go-webgpu/wgpu"
	wgpuext_glfw "github.com/rajveermalviya/go-webgpu/wgpuext/glfw"
)

func (w GLFWWindow) CreateSurface(instance *wgpu.Instance) *wgpu.Surface {
	return instance.CreateSurface(wgpuext_glfw.GetSurfaceDescriptor(w.handle))
}
