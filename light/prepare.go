package light

import "github.com/rajveermalviya/go-webgpu/wgpu"

func (l *Light) Prepare(groupIndex uint32, queue *wgpu.Queue, renderPass *wgpu.RenderPassEncoder) {
	l.UpdateUniform()
	l.WriteBuffer(queue)
	l.SetBindGroup(groupIndex, renderPass)
}
