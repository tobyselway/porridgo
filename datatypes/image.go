package datatypes

import (
	"fmt"
	"image/png"
	"os"
)

type Image struct {
	Pixels []byte
	Width  uint32
	Height uint32
}

func ImageFromPNG(filename string) (*Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening image file: %w", err)
	}
	defer file.Close()

	// Decode the PNG image
	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decoding PNG: %w", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	texels := make([]byte, width*4*height)

	// Populate pixel array
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			texels[y*4*width+x*4+0] = uint8(r >> 8)
			texels[y*4*width+x*4+1] = uint8(g >> 8)
			texels[y*4*width+x*4+2] = uint8(b >> 8)
			texels[y*4*width+x*4+3] = 0xff
		}
	}

	return &Image{
		Pixels: texels,
		Width:  uint32(width),
		Height: uint32(height),
	}, nil
}
