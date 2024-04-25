package texture

import (
	"porridgo/datatypes"
)

func FromPNG(filename string) (Texture, error) {
	image, err := datatypes.ImageFromPNG(filename)
	return Texture{
		image: image,
	}, err
}
