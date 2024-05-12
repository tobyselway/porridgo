package camera

import (
	"fmt"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (c *Camera) Setup(device *wgpu.Device, queue *wgpu.Queue) error {
	c.SetupUniform()
	err := c.UpdateUniform()
	if err != nil {
		return fmt.Errorf("setting up camera: %w", err)
	}

	err = c.SetupBuffer(device)
	if err != nil {
		return fmt.Errorf("setting up camera: %w", err)
	}

	err = c.CreateBindGroup(device)
	if err != nil {
		return fmt.Errorf("setting up camera: %w", err)
	}
	return nil
}
