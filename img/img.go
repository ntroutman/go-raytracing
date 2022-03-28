package img

import (
	"fmt"
	"os"
	"raytracing/vec"
)

type Image struct {
	Width  int
	Height int
	pixels []byte
}

func New(width, height int) Image {
	return Image{width, height, make([]byte, width*height*3)}
}

func Gradient(width, height int) Image {
	img := New(width, height)

	fwidth := float32(width - 1)
	fheight := float32(height - 1)
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			r := float32(x) / fwidth
			g := float32(y) / fheight
			b := float32(0.25)

			img.SetPixelRGB(x, y, r, g, b)
		}
	}

	return img
}

func (img *Image) SetPixelRGB(x, y int, r, g, b float32) {
	idx := img.Index(x, y)
	img.pixels[idx+2] = byte(256 * r)
	img.pixels[idx+1] = byte(256 * g)
	img.pixels[idx+0] = byte(256 * b)
}

func (img *Image) SetPixel(x, y int, color vec.Color) {
	img.SetPixelRGB(x, y, float32(color.X), float32(color.Y), float32(color.Z))
}

func (img *Image) Index(x, y int) int {
	return 3 * (y*img.Width + x)
}

func (img *Image) WriteTarga(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	// Format
	header := []byte{
		// 1 	1 byte 	ID length 	Length of the image ID field
		0,
		// 2 	1 byte 	Color map type 	Whether a color map is included
		0,
		// 3 	1 byte 	Image type 	Compression and color types
		2, // Raw image, no color map
		// 4 	5 bytes 	Color map specification 	Describes the color map
		0, 0, 0, 0, 0,
		// 5 	10 bytes 	Image specification 	Image dimensions and format
		// - 2 bytes    X origin
		0, 0,
		// - 2 bytes    Y origin
		0, 0,
		// - 2 bytes    Width
		(byte)(img.Width & 255), (byte)((img.Width >> 8) & 255),
		// - 2 bytes    Width
		(byte)(img.Height & 255), (byte)((img.Height >> 8) & 255),
		// - 1 byte     Pixel Depth (bits per Pixel)
		24,
		// - 1 byte     Image Descriptor
		0, //1 << 5, // Bit 5 = top-to-bottom ordering
	}
	file.Write(header)

	// 6 	From image ID length field 	Image ID 	Optional field containing identifying information
	// 7 	From color map specification field 	Color map data 	Look-up table containing color map data
	// 8 	From image specification field 	Image data 	Stored according to the image descriptor
	file.Write(img.pixels)
}
