package bitd_convert

import (
	"fmt"
	"os"

	"github.com/markhughes/dirry/internal/bitd"
	"github.com/markhughes/dirry/internal/palettes"
)

func Convert(filepath string, width int, height int, bitdepth int) {
	// read filepath into a bytes reader
	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	pallette, err := palettes.RetrievePallete(-1)
	if err != nil {
		panic(err)
	}

	info, bytes, err := bitd.ConvertImage(data, width, height, bitdepth, pallette)
	if err != nil {
		panic(err)
	}

	// save to file
	os.WriteFile(filepath+".png", bytes, 0644)

	fmt.Printf("Converted: %v\n", info)

}
