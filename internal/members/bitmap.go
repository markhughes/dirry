package members

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/markhughes/dirry/internal/palettes"
	"github.com/markhughes/dirry/internal/utils"
	"github.com/markhughes/dirry/internal/version"
)

type MemberBitmap struct {
	InitialRect  utils.Rect
	BoundingRect utils.Rect
	RegY         uint16
	RegX         uint16
	Flags        uint16
	Bytes        uint16
	BitsPerPixel uint16
	Pitch        uint16
	Clut         int16
	Alpha        uint8
	OLE          [7]byte

	CastMemberID int
}

func (m *MemberBitmap) ToJson() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (m *MemberBitmap) FromBytes(b []byte, v version.Version, flags uint8) error {
	var reader = bytes.NewReader(b)
	var err error

	utils.DebugMsg("members/bitmap", "\n-------------------------------------\n")
	utils.DebugMsg("members/bitmap", "flags: %v\n", flags)
	utils.DebugMsg("members/bitmap", "v: %v\n", v)

	if v.IsLessThan(version.Director_4_0_0) {

		m.Bytes, err = utils.ReadUInt16(reader, binary.BigEndian)
		if err != nil {
			return err
		}

		// TODO: why 0x7fff?
		m.Bytes = m.Bytes & 0x7fff

		m.InitialRect, err = utils.ReadRect(reader, binary.BigEndian)
		if err != nil {
			return err
		}

		m.BoundingRect, err = utils.ReadRect(reader, binary.BigEndian)
		if err != nil {
			return err
		}
		m.RegY, err = utils.ReadUInt16(reader, binary.BigEndian)
		if err != nil {
			return err
		}
		m.RegX, err = utils.ReadUInt16(reader, binary.BigEndian)
		if err != nil {
			return err
		}

		utils.DebugMsg("members/bitmap", "bytes: %v\n", m.Bytes)
		utils.DebugMsg("members/bitmap", "initialRect: %v\n", m.InitialRect)
		utils.DebugMsg("members/bitmap", "boundingRect: %v\n", m.BoundingRect)
		utils.DebugMsg("members/bitmap", "Y: %v\n", m.RegY)
		utils.DebugMsg("members/bitmap", "X: %v\n", m.RegX)

		if m.Bytes&0x8000 != 0 {

			m.BitsPerPixel, err = utils.ReadUInt16(reader, binary.BigEndian)
			if err != nil {
				return err
			}

			m.Clut, err = utils.ReadInt16(reader, binary.BigEndian)
			if err != nil {
				return err
			}

			if m.Clut <= 0 {
				// builtin palette ?
				m.Clut = m.Clut - 1
				m.CastMemberID = -1

			} else {
				m.CastMemberID = 1

			}
		} else {
			m.BitsPerPixel = 1
			m.Clut = int16(palettes.ClutSystemMac)
			m.CastMemberID = 1
		}

		var pitch = m.InitialRect.Width

		if pitch%16 != 0 {
			pitch += 16 - (m.InitialRect.Width % 16)
		}

		pitch = pitch * int16(m.BitsPerPixel)
		pitch >>= 3

		m.Pitch = uint16(pitch)

		// fmt.Printf("bitsPerPixel: %v\n", m.BitsPerPixel)
		// fmt.Printf("pitch: %v\n", pitch)

	} else if v.IsGreaterThanOrEqualTo(version.Director_4_0_0) && (v.IsLessThan(version.Director_6_0_0)) {
		binary.Read(reader, binary.BigEndian, &m.Bytes)
		m.Bytes &= 0x0fff

		initialRect, err := utils.ReadRect(reader, binary.BigEndian)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("initialRect: %v\n", initialRect)

		boundingRect, err := utils.ReadRect(reader, binary.BigEndian)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("boundingRect: %v\n", boundingRect)

		var Y uint16
		binary.Read(reader, binary.BigEndian, &Y)
		var X uint16
		binary.Read(reader, binary.BigEndian, &X)

		reader.ReadByte() // ?

		var bitsPerPixel uint8
		binary.Read(reader, binary.BigEndian, &bitsPerPixel)

		if reader.Len() > 0 {
			var clutCastLib int16
			if v.IsGreaterThanOrEqualTo(version.Director_5_0_0) {
				binary.Read(reader, binary.BigEndian, &clutCastLib)
			} else {
				clutCastLib = -1
			}

			binary.Read(reader, binary.BigEndian, &m.Clut)

			if m.Clut <= 0 {
				// built in palette ?
				m.Clut = m.Clut - 1
				m.CastMemberID = -1
			} else if m.Clut > 0 {
				if clutCastLib == -1 {
					// Leave nil, it will be populated when needed
					// m.CastMemberID = nil
				} else {
					m.CastMemberID = int(clutCastLib)
				}
			}
			var unk1, unk2, unk3 uint16
			var unk4, unk5 uint32

			binary.Read(reader, binary.BigEndian, &unk1)
			binary.Read(reader, binary.BigEndian, &unk2)
			binary.Read(reader, binary.BigEndian, &unk3)
			binary.Read(reader, binary.BigEndian, &unk4)
			binary.Read(reader, binary.BigEndian, &unk5)

			var flags uint16
			binary.Read(reader, binary.BigEndian, &flags)

			m.Flags = flags

		}

		if bitsPerPixel == 0 {
			bitsPerPixel = 1
		}

		var tail uint32
		var buf [256]byte = [256]byte{}
		for reader.Len() > 0 {
			var c byte
			binary.Read(reader, binary.BigEndian, &c)
			if tail < 256 {
				buf[tail] = c
			}
			tail++
		}

		m.InitialRect = initialRect
		m.BoundingRect = boundingRect
		m.RegY = Y
		m.RegX = X
		m.BitsPerPixel = uint16(bitsPerPixel)

	} else {

		var stride uint16
		binary.Read(reader, binary.BigEndian, &stride)

		rectangle, err := utils.ReadRect(reader, binary.BigEndian)
		if err != nil {
			panic(err)
		}

		var alphaThreshold uint8
		binary.Read(reader, binary.BigEndian, &alphaThreshold)

		var ole [7]byte
		binary.Read(reader, binary.BigEndian, &ole)

		var X uint16
		binary.Read(reader, binary.BigEndian, &X)

		var Y uint16
		binary.Read(reader, binary.BigEndian, &Y)

		var flags uint8
		binary.Read(reader, binary.BigEndian, &flags)

		var bitDepth uint8
		var palette int16
		var castLib int16

		if (stride & 0x8000) != 0 {
			stride &= 0x3FFF

			binary.Read(reader, binary.BigEndian, &bitDepth)

			binary.Read(reader, binary.BigEndian, &palette)
			binary.Read(reader, binary.BigEndian, &castLib)

		}
		m.Bytes = stride
		m.InitialRect = rectangle
		m.BoundingRect = rectangle
		m.RegY = Y
		m.RegX = X
		m.Alpha = alphaThreshold
		m.OLE = ole
		m.Flags = uint16(flags)
		m.BitsPerPixel = uint16(bitDepth)
		m.Clut = palette
		m.CastMemberID = int(castLib)

		// fmt.Printf("stride: %v\n", stride)
		// fmt.Printf("rectangle: %v\n", rectangle)
		// fmt.Printf("alphaThreshold: %v\n", alphaThreshold)
		// fmt.Printf("ole: %v\n", ole)
		// fmt.Printf("X: %v\n", X)
		// fmt.Printf("Y: %v\n", Y)
		// fmt.Printf("flags: %v\n", flags)
		// fmt.Printf("bitDepth: %v\n", bitDepth)
		// fmt.Printf("palette: %v\n", palette)
		// fmt.Printf("castLib: %v\n", castLib)

		/*
		   int v27 = input.ReadUInt16();

		     TotalWidth = v27 & 0x7FFF; //TODO: what does that last bit even do.. some sneaky flag?
		     //DIRAPI checks if TotalWidth & 0x8000 == 0

		     Rectangle = input.ReadRect();
		     AlphaThreshold = input.ReadByte();
		     OLE = input.ReadBytes(7).ToArray();

		     RegistrationPoint = input.ReadPoint();

		     Flags = (BitmapFlags)input.ReadByte();

		     if (!input.IsDataAvailable) return;
		     BitDepth = input.ReadByte();

		     if (!input.IsDataAvailable) return;
		     Palette = input.ReadInt32();
		*/
	}

	return nil

}
