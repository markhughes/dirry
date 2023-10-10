package members

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/markhughes/dirry/internal/utils"
	"github.com/markhughes/dirry/internal/version"
)

type MemberXtra struct {
	TypeLength int32
	Type       string
}

func (m *MemberXtra) ToJson() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
func (m *MemberXtra) FromBytes(b []byte, v version.Version, flags uint8) error {
	var err error

	var reader = bytes.NewReader(b)

	err = binary.Read(reader, binary.BigEndian, &m.TypeLength)
	if err != nil {
		return err
	}

	m.Type, err = utils.ReadString(reader, int(m.TypeLength), false)
	if err != nil {
		return err
	}
	utils.DebugMsg("members/xtra", "Type: %v\n", m.Type)

	return nil

}
