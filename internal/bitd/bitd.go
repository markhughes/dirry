package bitd

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/markhughes/dirry/internal/palettes"
	"github.com/markhughes/dirry/internal/utils"
)

type BitmapInfo struct {
	IsPackBitsCompressed bool
	PackbitsCompression  int
	Depth                int
}

func ConvertImage(data []byte, width, height, depth int, palette palettes.PaletteValue) (BitmapInfo, []byte, error) {
	utils.DebugMsg("bitd", "width: %v, height: %v, bitdepth: %v", width, height, depth)

	var err error
	var converted []byte

	buf := bytes.NewBuffer([]byte{})

	var info BitmapInfo
	switch depth {
	case 1:
		info, converted, err = Convert1BitImage(data, width, height)
		if err != nil {
			return info, nil, fmt.Errorf("failed to convert 1 bit image: %s", err)
		}

	case 2:
		info, converted, err = Convert2BitImage(data, width, height)
		if err != nil {
			return info, nil, fmt.Errorf("failed to convert 2 bit image: %s", err)
		}

	case 4:
		info, converted, err = Convert4BitImage(data, width, height, palette)
		if err != nil {
			return info, nil, fmt.Errorf("failed to convert 4 bit image: %s", err)
		}

	case 8:
		info, converted, err = Convert8BitImage(data, width, height, palette)
		if err != nil {
			return info, nil, fmt.Errorf("failed to convert 8 bit image: %s", err)
		}

	case 24:
		info, converted, err = Convert24BitImage(data, width, height)
		if err != nil {
			return info, nil, fmt.Errorf("failed to convert 24 bit image: %s", err)
		}

	case 32:
		info, converted, err = Convert32BitImage(data, width, height)
		if err != nil {
			return info, nil, fmt.Errorf("failed to convert 32 bit image: %s", err)
		}

	default:
		return BitmapInfo{}, nil, fmt.Errorf("unsupported bit depth: %d", depth)

	}

	if converted != nil {
		_, err = buf.Write(converted)
		if err != nil {
			return info, nil, fmt.Errorf("failed to write converted bit image to buffer: %s", err)
		}
	}

	return info, buf.Bytes(), nil
}

func Convert1BitImage(data []byte, width, height int) (BitmapInfo, []byte, error) {
	info := BitmapInfo{}
	info.Depth = 1

	width = ((width-1)/16 + 1) * 16 // 1-bit images width is multiple of 16

	uncompressedSize := (height * width * 1) / 8

	var err error
	if len(data) == uncompressedSize {
		data, err = unpackPackbits1(data)
		info.IsPackBitsCompressed = false
		info.PackbitsCompression = 1
	} else {
		data, err = unpackPackbits8(data)
		info.IsPackBitsCompressed = true
		info.PackbitsCompression = 8
	}

	if err != nil {
		return info, nil, fmt.Errorf("failed to decompress image data: %s", err)
	}

	img := image.NewGray(image.Rect(0, 0, width, height))

	// Iterate over the decompressed data..
	bitIndex := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			byteIndex := bitIndex / 8
			bitOffset := uint(bitIndex % 8)
			bit := data[byteIndex] & (0x80 >> bitOffset)
			if bit != 0 {
				img.Set(x, y, color.Gray{Y: 255})
			}

			bitIndex++
		}
	}

	var buf bytes.Buffer
	if png.Encode(&buf, img) != nil {
		return info, nil, fmt.Errorf("failed to encode image: %s", err)
	}

	return info, buf.Bytes(), nil
}

func Convert2BitImage(data []byte, width, height int) (BitmapInfo, []byte, error) {
	info := BitmapInfo{}
	info.Depth = 2

	uncompressedSize := (height * width * 2) / 8

	var err error
	if len(data) == uncompressedSize {
		data, err = unpackPackbits1(data)
		info.IsPackBitsCompressed = false
		info.PackbitsCompression = 1
	} else {
		data, err = unpackPackbits8(data)
		info.IsPackBitsCompressed = true
		info.PackbitsCompression = 8
	}

	if err != nil {
		return info, nil, fmt.Errorf("failed to decompress image data: %s", err)
	}

	img := image.NewGray(image.Rect(0, 0, width, height))

	// iterate over decompressed data
	bitIndex := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			byteIndex := bitIndex / 4
			bitOffset := uint(bitIndex % 4)
			bit := (data[byteIndex] & (0xC0 >> (bitOffset * 2))) >> ((3 - bitOffset) * 2)
			img.Set(x, y, color.Gray{Y: bit * 85}) // 85 = 255 / 3 to map 2-bit color to 8-bit color

			bitIndex++
		}
	}

	var buf bytes.Buffer
	if png.Encode(&buf, img) != nil {
		return info, nil, fmt.Errorf("failed to encode image: %s", err)
	}

	return info, buf.Bytes(), nil
}

func Convert4BitImage(data []byte, width, height int, palette palettes.PaletteValue) (BitmapInfo, []byte, error) {
	info := BitmapInfo{}
	info.Depth = 4

	width = ((width-1)/4 + 1) * 4

	uncompressedSize := (height * width * 4) / 8

	var err error
	if len(data) == uncompressedSize {
		data, err = unpackPackbits1(data)
		info.IsPackBitsCompressed = false
		info.PackbitsCompression = 1
	} else {
		data, err = unpackPackbits8(data)
		info.IsPackBitsCompressed = true
		info.PackbitsCompression = 8
	}

	if err != nil {
		return info, nil, fmt.Errorf("failed to decompress image data: %s", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Iterate over decompressed data
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the color index
			index := data[(y*width+x)/2]
			if (x % 2) == 0 {
				index >>= 4
			} else {
				index &= 0x0F
			}

			// Get the RGB values from our palette and set the pixel color in the image
			r := palette.Palette[index].R
			g := palette.Palette[index].G
			b := palette.Palette[index].B

			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	var buf bytes.Buffer
	if png.Encode(&buf, img) != nil {
		return info, nil, fmt.Errorf("failed to encode image: %s", err)
	}

	return info, buf.Bytes(), nil
}

func Convert8BitImage(data []byte, width, height int, palette palettes.PaletteValue) (BitmapInfo, []byte, error) {
	info := BitmapInfo{}
	info.Depth = 8

	uncompressedSize := (height * width * 8) / 8

	var err error
	if len(data) == uncompressedSize {
		data, err = unpackPackbits1(data)
		info.IsPackBitsCompressed = false
		info.PackbitsCompression = 1
	} else {
		data, err = unpackPackbits8(data)
		info.IsPackBitsCompressed = true
		info.PackbitsCompression = 8
	}

	if err != nil {
		return info, nil, fmt.Errorf("failed to decompress image data: %s", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// iterate over decompressed data
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the color index
			index := data[y*width+x]

			// Get the RGB values from palette
			r := palette.Palette[index].R
			g := palette.Palette[index].G
			b := palette.Palette[index].B

			// Set the pixel in the image
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	var buf bytes.Buffer
	if png.Encode(&buf, img) != nil {
		return info, nil, fmt.Errorf("failed to encode image: %s", err)
	}

	return info, buf.Bytes(), nil
}

func Convert24BitImage(data []byte, width, height int) (BitmapInfo, []byte, error) {
	info := BitmapInfo{}
	info.Depth = 24

	uncompressedSize := (height * width * 24) / 8

	var err error
	if len(data) == uncompressedSize {
		data, err = unpackPackbits1(data)
		info.IsPackBitsCompressed = false
		info.PackbitsCompression = 1
	} else {
		data, err = unpackPackbits24(data, width)
		info.IsPackBitsCompressed = true
		info.PackbitsCompression = 24
	}

	if err != nil {
		return info, nil, fmt.Errorf("failed to decompress image data: %s", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Iterate over decompressed data
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := data[(y*width+x)*3]
			g := data[(y*width+x)*3+1]
			b := data[(y*width+x)*3+2]

			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	var buf bytes.Buffer
	if png.Encode(&buf, img) != nil {
		return info, nil, fmt.Errorf("failed to encode image: %s", err)
	}

	return info, buf.Bytes(), nil
}

func Convert32BitImage(data []byte, width, height int) (BitmapInfo, []byte, error) {
	info := BitmapInfo{}
	info.Depth = 32

	uncompressedSize := (height * width * 32) / 8

	var err error
	if len(data) == uncompressedSize {
		data, err = unpackPackbits2(data)
		info.IsPackBitsCompressed = false
		info.PackbitsCompression = 2
	} else {
		data, err = unpackPackbits32(data, width)
		info.IsPackBitsCompressed = true
		info.PackbitsCompression = 32
	}

	if err != nil {
		return info, nil, fmt.Errorf("failed to decompress image data: %s", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Iterate over the decompressed stuff
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := data[(y*width+x)*4+1]
			g := data[(y*width+x)*4+2]
			b := data[(y*width+x)*4+3]
			a := 255 - data[(y*width+x)*4]

			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	var buf bytes.Buffer
	if png.Encode(&buf, img) != nil {
		return info, nil, fmt.Errorf("failed to encode image: %s", err)
	}

	return info, buf.Bytes(), nil
}
