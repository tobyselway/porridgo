package texture

func (t Texture) LabelID() string {
	switch t.textureType {
	case Image:
		return t.image.Filename
	case Depth:
		return "Depth"
	}
	return "Unknown"
}
