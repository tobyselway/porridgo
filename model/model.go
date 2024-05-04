package model

import (
	"porridgo/material"
	"porridgo/mesh"
)

type Model struct {
	Name      string
	Meshes    []mesh.Mesh
	Materials []material.Material
}

func (m *Model) Cleanup() {
	for i := range m.Materials {
		defer m.Materials[i].Cleanup()
	}
	for i := range m.Meshes {
		defer m.Meshes[i].Cleanup()
	}
}
