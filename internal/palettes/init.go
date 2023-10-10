//go:build !js

package palettes

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/markhughes/dirry/internal/consts"
)

func init() {

	palSystemMac, err := fromFile(path.Join(consts.PalettesDir, "SystemMac.json"))
	if err != nil {
		panic(err)
	}
	RegisterPallete(ClutSystemMac, palSystemMac)

	palRainbow, err := fromFile(path.Join(consts.PalettesDir, "Rainbow.json"))
	if err != nil {
		panic(err)
	}
	RegisterPallete(ClutRainbow, palRainbow)

	palGrayscale, err := fromFile(path.Join(consts.PalettesDir, "Grayscale.json"))
	if err != nil {
		panic(err)
	}
	RegisterPallete(ClutGrayscale, palGrayscale)

	palPastels, err := fromFile(path.Join(consts.PalettesDir, "Pastels.json"))
	if err != nil {
		panic(err)
	}
	RegisterPallete(ClutPastels, palPastels)

	palVivid, err := fromFile(path.Join(consts.PalettesDir, "Vivid.json"))
	if err != nil {
		panic(err)
	}
	RegisterPallete(ClutVivid, palVivid)

	palNTSC, err := fromFile(path.Join(consts.PalettesDir, "NTSC.json"))
	if err != nil {
		panic(err)
	}
	RegisterPallete(ClutNTSC, palNTSC)

	palMetallic, err := fromFile(path.Join(consts.PalettesDir, "Metallic.json"))
	if err != nil {
		panic(err)
	}
	RegisterPallete(ClutMetallic, palMetallic)

	palSystemWin, err := fromFile(path.Join(consts.PalettesDir, "SystemWin.json"))
	if err != nil {
		panic(err)
	}
	RegisterPallete(ClutSystemWin, palSystemWin)

	palSystemWinD5, err := fromFile(path.Join(consts.PalettesDir, "SystemWinD5.json"))
	if err != nil {
		panic(err)
	}
	RegisterPallete(ClutSystemWinD5, palSystemWinD5)

}

func fromFile(path string) (PaletteValue, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return PaletteValue{}, err
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		return PaletteValue{}, err
	}

	// Define a slice of PalettePixel to hold the unmarshalled JSON
	var palettePixels []Pixel24
	err = json.Unmarshal(bytes, &palettePixels)
	if err != nil {
		return PaletteValue{}, err
	}

	// Convert the slice of PalettePixel to a PaletteValue
	var paletteValue PaletteValue
	paletteValue.Size = int32(len(palettePixels))
	for i, pixel := range palettePixels {
		paletteValue.Palette[i] = Pixel24{
			R: pixel.R,
			G: pixel.G,
			B: pixel.B,
		}
	}

	return paletteValue, nil
}
