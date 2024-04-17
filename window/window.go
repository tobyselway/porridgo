package window

import "github.com/rajveermalviya/go-webgpu/wgpu"

type Window interface {
	Open() error

	Destroy()

	OnResize(callback func(int, int))

	OnKeyEvent(callback func(Key, int, Action, ModifierKey))

	Size() (int, int)

	Run(tickCallback func()) error

	CreateSurface(instance *wgpu.Instance) *wgpu.Surface
}
