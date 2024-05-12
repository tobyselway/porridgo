package light

import (
	"fmt"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (l *Light) Setup(device *wgpu.Device, queue *wgpu.Queue) error {
	l.SetupUniform()
	err := l.UpdateUniform()
	if err != nil {
		return fmt.Errorf("setting up light: %w", err)
	}

	err = l.SetupBuffer(device)
	if err != nil {
		return fmt.Errorf("setting up light: %w", err)
	}

	err = l.CreateBindGroup(device)
	if err != nil {
		return fmt.Errorf("setting up light: %w", err)
	}
	return nil
}
