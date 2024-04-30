package vertex

import "porridgo/datatypes"

type Vertex struct {
	Position  datatypes.Vec3f
	UVMapping datatypes.Vec2f
}

func NewVertex(posX float32, posY float32, posZ float32, uvX float32, uvY float32) Vertex {
	return Vertex{
		Position:  datatypes.NewVec3f(posX, posY, posZ),
		UVMapping: datatypes.NewVec2f(uvX, uvY),
	}
}
