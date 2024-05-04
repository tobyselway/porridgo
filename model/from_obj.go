package model

import (
	"fmt"
	"path/filepath"
	"porridgo/datatypes"
	"porridgo/material"
	"porridgo/mesh"
	"porridgo/texture"

	"github.com/udhos/gwob"
)

func FromOBJ(filename string) (Model, error) {
	name := filename
	filename, err := filepath.Abs(filename)
	if err != nil {
		return Model{}, fmt.Errorf("determining absolute path: %w", err)
	}

	obj, err := gwob.NewObjFromFile(filename, &gwob.ObjParserOptions{})
	if err != nil {
		return Model{}, fmt.Errorf("loading obj: %w", err)
	}

	mtlFilename := filepath.Join(filepath.Dir(filename), obj.Mtllib)

	mtl, err := gwob.ReadMaterialLibFromFile(mtlFilename, &gwob.ObjParserOptions{})
	if err != nil {
		return Model{}, fmt.Errorf("loading obj mtl: %w", err)
	}

	materials := []material.Material{}

	materialsNameToIdx := map[string]uint32{}

	for name, mtl := range mtl.Lib {
		diffuseTextureFilename := filepath.Join(filepath.Dir(filename), mtl.MapKd)
		diffuseTexture, err := texture.FromFile(diffuseTextureFilename)
		if err != nil {
			return Model{}, fmt.Errorf("loading obj mtl diffuse map: %w", err)
		}

		materials = append(materials, material.Material{
			Name:           name,
			DiffuseTexture: diffuseTexture,
		})
		materialsNameToIdx[name] = uint32(len(materials)) - 1
	}

	meshes := []mesh.Mesh{}
	for _, group := range obj.Groups {
		vertices := []mesh.Vertex{}

		for i := 0; i < len(obj.Coord)/8; i++ {
			vertices = append(vertices, mesh.Vertex{
				Position: datatypes.NewVec3f(
					obj.Coord[i*8+0],
					obj.Coord[i*8+1],
					obj.Coord[i*8+2],
				),
				TexCoords: datatypes.NewVec2f(
					obj.Coord[i*8+3],
					1.0-obj.Coord[i*8+4],
				),
				Normal: datatypes.NewVec3f(
					obj.Coord[i*8+5],
					obj.Coord[i*8+6],
					obj.Coord[i*8+7],
				),
			})
		}

		original := obj.Indices[group.IndexBegin:(group.IndexBegin + group.IndexCount)]
		indices := []uint32{}

		for _, index := range original {
			indices = append(indices, uint32(index))
		}

		meshes = append(meshes, mesh.Mesh{
			Name:        group.Name,
			NumElements: uint32(group.IndexCount),
			Material:    materialsNameToIdx[group.Usemtl],
			Vertices:    vertices,
			Indices:     indices,
		})
	}

	return Model{
		Name:      name,
		Materials: materials,
		Meshes:    meshes,
	}, nil
}
