package shader

import (
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type Shader struct {
	Module *wgpu.ShaderModule
	Name   string
	Code   string
}

func (s *Shader) Cleanup() {
	if s.Module != nil {
		defer s.Module.Release()
		s.Module = nil
	}
}
