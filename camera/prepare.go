package camera

import "github.com/rajveermalviya/go-webgpu/wgpu"

func (c *Camera) Prepare(groupIndex uint32, queue *wgpu.Queue, renderPass *wgpu.RenderPassEncoder) {
	c.UpdateUniform()
	c.WriteBuffer(queue)
	c.SetBindGroup(groupIndex, renderPass)
}
