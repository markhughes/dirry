package chunks

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
)

type XtrlChunk struct {
	Chunk  *BinaryChunk
	Length int32
	List   XtraList
}

func ReadXtrlChunkRaw(r *bytes.Reader, endian binary.ByteOrder) (*XtrlChunk, error) {

	var err error
	chunk := &XtrlChunk{}
	chunk.Chunk, err = FromBinaryAtPartHeadless(r, "XTRl", endian)
	if err != nil {
		return nil, err
	}

	err = processXtrlChunk(chunk, binary.BigEndian)
	if err != nil {
		return nil, err
	}

	return chunk, nil

}

// TODO: I broke this
func processXtrlChunk(chunk *XtrlChunk, endian binary.ByteOrder) error {
	var err error

	chunk.Chunk.ReadBytes(4)

	var xtra XtraList
	xtra.XtraCount, err = chunk.Chunk.ReadUInt32(endian)
	if err != nil {
		return err
	}

	// fmt.Printf("xtra count: %d\n", xtra.XtraCount)

	// xtra.Entries = make([]Entry, xtra.XtraCount)

	chunk.List = xtra
	return nil

	// for i := 0; i < int(xtra.XtraCount); i++ {
	// 	// currentPos := chunk.Chunk.Position()

	// 	offset, err := chunk.Chunk.ReadUInt32(endian)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	fmt.Printf("offset: %d, position: %d\n", offset, chunk.Chunk.Position())

	// 	// length := int64(offset) - currentPos

	// 	q1, _ := chunk.Chunk.ReadInt32(endian) // ?
	// 	q2, _ := chunk.Chunk.ReadInt32(endian) // ?

	// 	fmt.Printf("q1: %d, q2: %d\n", q1, q2)

	// 	xtra.Entries[i].Guid, err = chunk.Chunk.ReadBytes(16)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	offsetCount, err := chunk.Chunk.ReadInt16(endian)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	offsets := make([]int32, offsetCount)
	// 	for j := 0; j < int(offsetCount); j++ {
	// 		offsets[j], err = chunk.Chunk.ReadInt32(endian)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}

	// 	// fmt.Printf("offsets: %+v\n", offsets)

	// 	var i = 0
	// 	for chunk.Chunk.Position() < int64(offset) {
	// 		xtra.Entries[i].Values = make([]VList16, 0)
	// 		xtra.Entries[i].Values[i].Unknown, err = chunk.Chunk.ReadByte()
	// 		if err != nil {
	// 			break
	// 		}

	// 		xtra.Entries[i].Values[i].Kind, err = chunk.Chunk.ReadByte()
	// 		if err != nil {
	// 			return nil, fmt.Errorf("could not read kind: %v", err)
	// 		}

	// 		len, err := chunk.Chunk.ReadVarInt()
	// 		if err != nil {
	// 			return nil, fmt.Errorf("could not read length: %v", err)
	// 		}

	// 		xtra.Entries[i].Values[i].Value, err = chunk.Chunk.ReadString(int(len), endian)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("could not read value: %v", err)
	// 		}

	// 		chunk.Chunk.ReadByte()

	// 		i++
	// 	}
	// }

	// if err != nil {
	// 	fmt.Println("binary.Read failed:", err)
	// }

	// fmt.Printf("XtraList: %+v\n", xtra)

}

func ReadXtrlChunk(r *os.File, endian binary.ByteOrder, offset int64) (*XtrlChunk, error) {
	_, err := r.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	chunk := &XtrlChunk{}

	chunk.Chunk, err = FromBinaryAt(r, "XTRl", endian, offset)
	if err != nil {
		return nil, err
	}

	err = processXtrlChunk(chunk, endian)
	if err != nil {
		return nil, err
	}

	return chunk, nil

}

func (c *XtrlChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

/*
This is a series of buffers with a single Vector List each, describing the different Xtras for the movie. This draws some information from xtrainfo.txt but additionally has the GUID of each required component. Corresponds to Lingo's the movieXtraList property.
skip(4)
INT32 vListsCount = xtraListChunkBuffer.readInt32();
for(var i=0;i<vListsCount;i++) {
INT32 vListCount = xtraListChunkBuffer.readInt32();
VList16 vList = xtraListChunkBuffer.readVList16();
Buffer theGUID = vList.numbers[2] + vList.numbers[3] + vList.numbers[4] + vList.numbers[5] + vList.numbers[6] + vList.numbers[7] + vList.numbers[8] + vList.numbers[9] + vList.numbers[10]
for(var j=0;j<vList.length;j++) {
vList[j].skip(2)
UINT8 xtraInfoNameLength
	String xtraInfoName = vList[j].readString(xtraInfoNameLength)
}
}

*/

type XtraList struct {
	Unknown   uint32
	XtraCount uint32
	Entries   []Entry
}

type Entry struct {
	Guid   []byte
	Values []VList16
}

type VList16 struct {
	Unknown byte
	Kind    byte
	Value   string
}
