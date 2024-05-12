package shader

import "github.com/rajveermalviya/go-webgpu/wgpu"

func (s *Shader) Setup(device *wgpu.Device) error {
	var err error
	s.Module, err = device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label:          s.Name,
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: s.Code},
	})
	if err != nil {
		return err
	}
	return nil
}
