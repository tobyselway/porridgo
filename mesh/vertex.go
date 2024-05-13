package mesh

import (
	"fmt"
	"porridgo/datatypes"
	"strings"
)

type Vertex struct {
	Position  datatypes.Vec3f
	TexCoords datatypes.Vec2f
	Normal    datatypes.Vec3f
	Tangent   datatypes.Vec3f
	Bitangent datatypes.Vec3f
}

func (v Vertex) String() string {
	return fmt.Sprintf("{ pos: %v, tex: %v, normal: %v }", v.Position, v.TexCoords, v.Normal)
}

type VertexArray []Vertex

func (va VertexArray) String() string {
	elems := []string{""}
	for _, v := range va {
		elems = append(elems, fmt.Sprintf("- %v", v))
	}
	return strings.Join(elems, "\n")
}
