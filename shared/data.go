package shared

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)
type Dir int32 

const (
	Up Dir = iota
	Down 
	Left
	Right 
)

type Player struct {
	SessionID  string
	XPosititon int
	YPosition  int
}
type PacketType rune

const (
	PacketTypeConnect    PacketType = 'C'
	PacketTypeMap        PacketType = 'M'
	PacketTypeMove       PacketType = 'L'
	PacketTypeDisconnect PacketType = 'X'
	PacketTypeErr        PacketType = 'E'
	PacketTypeSuccess	PacketType = 'S'
)

type Cell int32 // SEE BELOW vvvvv
const (
	Empty Cell = iota // 0
	Land // 1
	Water // 2 
	Mountains // 3 
	User // 4

)
// 0 for empty
// 1 for land
// 2 for water
// 3 for mountains
// 4 for player

type Universe struct {
	Map  []Cell
	Size int // map width = height
}

func NewPlayer() *Player {
	return &Player{XPosititon: 0, YPosition: 0}
}

func ConvertMapToBytes(u Universe) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, cell := range u.Map {
		if err := binary.Write(buf, binary.LittleEndian, cell); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func CompressMapData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	writer.Close()
	return buf.Bytes(), nil
}

func DecompressMapData(data []byte) ([]byte, error) {
    // fmt.Printf("%s\n",data) // Log data length
    
    buf := bytes.NewReader(data)
    reader, err := zlib.NewReader(buf)
    if err != nil {
        fmt.Printf("creating newreader, %s", err)
        return nil, err
    }
    defer reader.Close()

    var out bytes.Buffer
    _, err = io.Copy(&out, reader)
    if err != nil {
		// errStr := fmt.Errorf("ERROR copying %s\nof len %d, %s", data, len(data), err)
		errStr := fmt.Errorf("ERROR")
        return nil, errStr
    }

    decompressedData := out.Bytes()
    
    return decompressedData, nil
}

func ConvertBytesToMap(mapSize int, data []byte) (Universe, error) {

	u := Universe{
		Map:  make([]Cell, mapSize*mapSize),
		Size: mapSize,
	}
	buf := bytes.NewReader(data)
	for i := range u.Map {
		if err := binary.Read(buf, binary.LittleEndian, &u.Map[i]); err != nil {
			return u, err
		}
	}
	return u, nil
}
