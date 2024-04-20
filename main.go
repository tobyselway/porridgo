package main

import (
	"fmt"
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

	r, err := renderer.CreateRenderer(&w, renderer.Config{})
	if err != nil {
		panic(err)
	}
	defer r.Cleanup()

	w.OnKeyEvent(func(key window.Key, scancode int, action window.Action, modifier window.ModifierKey) {
		// Print resource usage on pressing 'R'
		if key == window.KeyR && (action == window.Press || action == window.Repeat) {
			r.PrintReport()
		}
	})

	w.Run(func() {
		err := r.Render()
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
