package camera

import (
	"porridgo/label"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (c *Camera) SetupBuffer(device *wgpu.Device) error {
	var err error
	c.buffer, err = device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    label.Label(c, "Buffer"),
		Contents: wgpu.ToBytes([]Uniform{c.uniform}),
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	return err
}

func (c *Camera) WriteBuffer(queue *wgpu.Queue) error {
	return queue.WriteBuffer(c.buffer, 0, wgpu.ToBytes([]Uniform{c.uniform}))
}
