package palettes

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/markhughes/dirry/internal/consts"
)

type Clut int16

type Pixel24 struct {
	R uint8 `json:"R"`
	G uint8 `json:"G"`
	B uint8 `json:"B"`
}

type PaletteValue struct {
	Size    int32
	Palette [256]Pixel24
}

func (pal *PaletteValue) ToColorTable() []byte {
	colorTable := make([]byte, 0, 4*256)

	for _, pixel := range pal.Palette {
		colorTable = append(colorTable, pixel.B, pixel.G, pixel.R, 0x00)
	}

	return colorTable
}

func (pal *PaletteValue) ToHtmlDoc() []byte {
	var out = bytes.NewBufferString("<html><body><table>")
	for i, pixel := range pal.Palette {
		out.WriteString(fmt.Sprintf("<tr><td>%d</td><td style=\"background-color: rgb(%d, %d, %d);\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</td></tr>\n", i, pixel.R, pixel.G, pixel.B))
	}

	out.WriteString("</table></body></html>")

	return out.Bytes()
}

const (
	ClutSystemMac   Clut = -1
	ClutRainbow     Clut = -2
	ClutGrayscale   Clut = -3
	ClutPastels     Clut = -4
	ClutVivid       Clut = -5
	ClutNTSC        Clut = -6
	ClutMetallic    Clut = -7
	ClutSystemWin   Clut = -101
	ClutSystemWinD5 Clut = -102
)

var registerdPalettes = map[Clut]PaletteValue{}

func RegisterPallete(clut Clut, pal PaletteValue) {
	if pal.Size == 0 {
		pal.Size = int32(len(pal.Palette))
	}
	registerdPalettes[clut] = pal

	DumpPalleteDebug(pal, clut)
}

func DumpPalleteDebug(pal PaletteValue, clut Clut) {
	// Store as a HTML doc for reference
	var out = bytes.NewBufferString("<html><body><table>")
	out.Write(pal.ToHtmlDoc())
	out.WriteString("</table></body></html>")

	os.Mkdir(path.Join(consts.PathDump, "_debug"), 0755)

	filepath := path.Join(consts.PathDump, "_debug", "palette_"+fmt.Sprintf("%d", clut))
	os.WriteFile(filepath+".html", out.Bytes(), 0644)

}

func RetrievePallete(clut Clut) (PaletteValue, error) {
	pallete, ok := registerdPalettes[clut]
	if !ok {
		return PaletteValue{}, fmt.Errorf("clut %d not found", clut)
	}
	return pallete, nil

}

func StoreAsHtml() {

}
