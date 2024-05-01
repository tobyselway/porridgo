package texture

import (
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (t *Texture) Setup(device *wgpu.Device, queue *wgpu.Queue) error {
	switch t.textureType {
	case Image:
		return t.setupImage(device, queue)
	case Depth:
		return t.setupDepth(device)
	}
	return nil
}
