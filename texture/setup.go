package texture

import (
	"fmt"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (t *Texture) Setup(device *wgpu.Device, queue *wgpu.Queue) error {
	var err error
	switch t.textureType {
	case Image:
		err = t.setupImage(device, queue)
		if err != nil {
			return fmt.Errorf("setting up texture %s: %w", t.image.Filename, err)
		}
		err = t.CreateBindGroup(device)
	case Depth:
		err = t.setupDepth(device)
	}
	if err != nil {
		return fmt.Errorf("setting up depth texture: %w", err)
	}

	return nil
}
