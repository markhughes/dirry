package shockwave

import "encoding/json"

type CodecType int

const (
	Standard CodecType = iota
	Afterburner
)

func (et CodecType) String() string {
	return [...]string{"Standard", "Afterburner"}[et]
}

type Codec struct {
	Name string
	Type CodecType
}

func (c Codec) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil

}

var codecsMap = map[string]Codec{
	"MV85": {
		Name: "MV85",
		Type: Standard,
	},

	// Director 6, Director 6.5?
	"MV93": {
		Name: "MV93",
		Type: Standard,
	},
	"MC93": {
		Name: "MC93",
		Type: Standard,
	},
	"MC95": {
		Name: "MC95",
		Type: Standard,
	},
	"FGDM": {
		Name: "FGDM",
		Type: Afterburner,
	},
	"FGDC": {
		Name: "FGDC",
		Type: Afterburner,
	},
}

func CodecByName(name string) (Codec, bool) {
	codec, ok := codecsMap[name]
	return codec, ok
}
