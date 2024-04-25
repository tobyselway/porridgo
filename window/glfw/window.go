package glfw

import (
	"fmt"
	"math"

	"porridgo/window"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type GLFWWindow struct {
	handle *glfw.Window
	Title  string
	Width  int
	Height int
}

func (w *GLFWWindow) Open() error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("initializing glfw: %w", err)
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	window, err := glfw.CreateWindow(w.Width, w.Height, w.Title, nil, nil)
	if err != nil {
		return fmt.Errorf("creating glfw window: %w", err)
	}
	w.handle = window
	return nil
}

func (w *GLFWWindow) Cleanup() {
	// Defenestrate
	w.handle.Destroy()
	glfw.Terminate()
}

func (w *GLFWWindow) OnResize(callback func(int, int)) {
	w.handle.SetSizeCallback(func(_ *glfw.Window, width int, height int) {
		callback(width, height)
	})
}

func (w *GLFWWindow) OnKeyEvent(callback func(window.Key, int, window.Action, window.ModifierKey)) {
	w.handle.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		callback(keyMapping[key], scancode, actionMapping[action], modifierKeyMapping[mods])
	})
}

func (w GLFWWindow) Cursor() (int, int) {
	x, y := w.handle.GetCursorPos()
	return int(math.Floor(x)), int(math.Floor(y))
}

func (w GLFWWindow) Size() (int, int) {
	return w.handle.GetSize()
}

func (w GLFWWindow) Run(tickCallback func()) error {
	for !w.handle.ShouldClose() {
		glfw.PollEvents()
		tickCallback()
	}
	return nil
}
