package light

import (
	"porridgo/label"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func (l *Light) SetupBuffer(device *wgpu.Device) error {
	var err error
	l.buffer, err = device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    label.Label(l, "Buffer"),
		Contents: wgpu.ToBytes([]Uniform{l.uniform}),
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	return err
}

func (l *Light) WriteBuffer(queue *wgpu.Queue) error {
	return queue.WriteBuffer(l.buffer, 0, wgpu.ToBytes([]Uniform{l.uniform}))
}
