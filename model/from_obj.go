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
		return Model{}, fmt.Errorf("determining absolute file path for obj %s: %w", name, err)
	}

	obj, err := gwob.NewObjFromFile(filename, &gwob.ObjParserOptions{})
	if err != nil {
		return Model{}, fmt.Errorf("loading obj %s: %w", name, err)
	}

	materials := []material.Material{}
	materialsNameToIdx := map[string]uint32{}

	if obj.Mtllib != "" {
		mtlFilename := filepath.Join(filepath.Dir(filename), obj.Mtllib)

		mtl, err := gwob.ReadMaterialLibFromFile(mtlFilename, &gwob.ObjParserOptions{})
		if err != nil {
			return Model{}, fmt.Errorf("loading obj %s: mtl file: %w", name, err)
		}

		for mtlName, mtl := range mtl.Lib {
			var diffuseTexture *texture.Texture = nil
			if mtl.MapKd != "" {
				diffuseTextureFilename := filepath.Join(filepath.Dir(filename), mtl.MapKd)
				dt, err := texture.FromFile(diffuseTextureFilename)
				diffuseTexture = &dt
				if err != nil {
					return Model{}, fmt.Errorf("loading obj %s: mtl %s: diffuse map: %w", name, mtlName, err)
				}
			}

			var normalTexture *texture.Texture = nil
			if mtl.Bump != "" {
				normalTextureFilename := filepath.Join(filepath.Dir(filename), mtl.Bump)
				dt, err := texture.FromFile(normalTextureFilename)
				normalTexture = &dt
				if err != nil {
					return Model{}, fmt.Errorf("loading obj %s: mtl %s: normal map: %w", name, mtlName, err)
				}
			}

			materials = append(materials, material.Material{
				Name:           mtlName,
				DiffuseTexture: diffuseTexture,
				NormalTexture:  normalTexture,
			})
			materialsNameToIdx[mtlName] = uint32(len(materials)) - 1
		}

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
				Tangent:   datatypes.NewVec3f(0, 0, 0),
				Bitangent: datatypes.NewVec3f(0, 0, 0),
			})
		}

		original := obj.Indices[group.IndexBegin:(group.IndexBegin + group.IndexCount)]
		indices := []uint32{}

		for _, index := range original {
			indices = append(indices, uint32(index))
		}

		var materialIdx *uint32 = nil

		idx, ok := materialsNameToIdx[group.Usemtl]
		if ok {
			materialIdx = &idx
		}

		vertices = calcTangents(vertices, indices)

		meshes = append(meshes, mesh.Mesh{
			Name:        group.Name,
			NumElements: uint32(group.IndexCount),
			Material:    materialIdx,
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

func calcTangents(vertices []mesh.Vertex, indices []uint32) []mesh.Vertex {
	trianglesIncluded := make([]uint32, len(vertices))

	// Calculate tangents and bitangets. We're going to
	// use the triangles, so we need to loop through the
	// indices in chunks of 3
	for i := 0; i < len(indices)/3; i++ {
		c0 := indices[i*3+0]
		c1 := indices[i*3+1]
		c2 := indices[i*3+2]

		v0 := vertices[c0]
		v1 := vertices[c1]
		v2 := vertices[c2]

		pos0 := v0.Position
		pos1 := v1.Position
		pos2 := v2.Position

		uv0 := v0.TexCoords
		uv1 := v1.TexCoords
		uv2 := v2.TexCoords

		// Calculate the edges of the triangle
		deltaPos1 := pos1.Sub(pos0)
		deltaPos2 := pos2.Sub(pos0)

		// This will give us a direction to calculate the
		// tangent and bitangent
		deltaUv1 := uv1.Sub(uv0)
		deltaUv2 := uv2.Sub(uv0)

		// Solving the following system of equations will
		// give us the tangent and bitangent.
		//     delta_pos1 = delta_uv1.x * T + delta_u.y * B
		//     delta_pos2 = delta_uv2.x * T + delta_uv2.y * B
		// Luckily, the place I found this equation provided
		// the solution!
		r := 1.0 / (deltaUv1.X*deltaUv2.Y - deltaUv1.Y*deltaUv2.X)
		tangent := (deltaPos1.MulScalar(deltaUv2.Y).Sub(deltaPos2.MulScalar(deltaUv1.Y))).MulScalar(r)
		// We flip the bitangent to enable right-handed normal
		// maps with wgpu texture coordinate system
		bitangent := (deltaPos2.MulScalar(deltaUv1.X).Sub(deltaPos1.MulScalar(deltaUv2.X))).MulScalar(-r)

		// We'll use the same tangent/bitangent for each vertex in the triangle
		vertices[c0].Tangent = tangent.Add(v0.Tangent)
		vertices[c1].Tangent = tangent.Add(v1.Tangent)
		vertices[c2].Tangent = tangent.Add(v2.Tangent)
		vertices[c0].Bitangent = bitangent.Add(v0.Bitangent)
		vertices[c1].Bitangent = bitangent.Add(v1.Bitangent)
		vertices[c2].Bitangent = bitangent.Add(v2.Bitangent)

		// Used to average the tangents/bitangents
		trianglesIncluded[c0] += 1
		trianglesIncluded[c1] += 1
		trianglesIncluded[c2] += 1
	}

	// Average the tangents/bitangents
	for i, n := range trianglesIncluded {
		denom := 1.0 / float32(n)
		vertices[i].Tangent = vertices[i].Tangent.MulScalar(denom)
		vertices[i].Bitangent = vertices[i].Bitangent.MulScalar(denom)
	}

	return vertices
}
