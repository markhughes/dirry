package members

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/markhughes/dirry/internal/utils"
	"github.com/markhughes/dirry/internal/version"
)

type MemberScript struct {
	Unknown1   uint8
	ScriptType uint8
}

func (m *MemberScript) ToJson() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
func (m *MemberScript) FromBytes(b []byte, v version.Version, flags uint8) error {
	// TODO: everyone is doing this differently it seems, so this is probably very wrong
	var err error
	var reader = bytes.NewReader(b)

	err = binary.Read(reader, binary.BigEndian, &m.Unknown1)
	if err != nil {
		return err
	}

	err = binary.Read(reader, binary.BigEndian, &m.ScriptType)
	if err != nil {
		return err
	}

	utils.DebugMsg("members/script", "Unknown1: %v\n", m.Unknown1)
	utils.DebugMsg("members/script", "ScriptType: %v\n", m.ScriptType)

	return nil

}
