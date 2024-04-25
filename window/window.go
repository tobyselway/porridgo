package window

import (
	"porridgo/lifecycle"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Window interface {
	lifecycle.RequiresCleanup

	Open() error

	OnResize(callback func(int, int))

	OnKeyEvent(callback func(Key, int, Action, ModifierKey))

	Size() (int, int)

	Run(tickCallback func()) error

	CreateSurface(instance *wgpu.Instance) *wgpu.Surface

	Cursor() (int, int)
}
