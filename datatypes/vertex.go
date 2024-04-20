package datatypes

type Vertex struct {
	Position  Vec3f
	UVMapping Vec2f
}

func NewVertex(posX float32, posY float32, posZ float32, uvX float32, uvY float32) Vertex {
	return Vertex{
		Position:  NewVec3f(posX, posY, posZ),
		UVMapping: NewVec2f(uvX, uvY),
	}
}
