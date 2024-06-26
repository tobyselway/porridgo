package main

import (
	"fmt"
	"porridgo/camera"
	"porridgo/datatypes"
	"porridgo/light"
	"porridgo/model"
	"porridgo/renderer"
	"porridgo/window"
	glfw_window "porridgo/window/glfw"
	"strings"

	_ "embed"
)

func main() {
	w := glfw_window.GLFWWindow{
		Title:  "porridgo",
		Width:  640,
		Height: 480,
	}
	err := w.Open()
	if err != nil {
		panic(err)
	}
	defer w.Cleanup()

	cam := camera.Camera{
		Position: datatypes.NewVec3f(0.0, 0.0, -2.0),
		Pitch:    0.0,
		Yaw:      0.0,
		Aspect:   float32(w.Width) / float32(w.Height),
		FovY:     45.0,
		ZNear:    0.1,
		ZFar:     100.0,
	}

	// mdl, err := model.FromOBJ("assets/rounded-cube.obj")
	mdl, err := model.FromOBJ("assets/cube.obj")
	if err != nil {
		panic(err)
	}

	sun := light.Light{
		Position: datatypes.NewVec3f(2.0, 2.0, 2.0),
		Color:    datatypes.NewVec3f(1.0, 1.0, 1.0),
	}

	r, err := renderer.CreateRenderer(&w, &cam, &sun, &mdl, renderer.Config{})
	if err != nil {
		panic(err)
	}
	defer r.Cleanup()

	camController := camera.Controller{
		Camera: &cam,
		Speed:  0.1,
	}

	spacePressed := false

	w.OnKeyEvent(func(key window.Key, scancode int, action window.Action, modifier window.ModifierKey) {
		// Print resource usage on pressing 'R'
		if key == window.KeyR && (action == window.Press || action == window.Repeat) {
			r.PrintReport()
		}

		if key == window.KeySpace && (action == window.Press || action == window.Repeat) {
			spacePressed = true
		} else {
			spacePressed = false
		}

		camController.ProcessKey(key, action, modifier)
	})

	w.Run(func() {
		cursorX, cursorY := w.Cursor()
		camController.ProcessMouse(float32(cursorX)/float32(w.Width), float32(cursorY)/float32(w.Height))
		camController.UpdateCamera()
		err := r.Render(spacePressed)
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
	})
}
