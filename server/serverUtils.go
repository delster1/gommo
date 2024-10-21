package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"gommo/shared"
	"net"
)

type Server struct {
	status      int
	Playerlist  map[string]*shared.Player
	port        int
	u           shared.Universe
	Connections []Connection
}
type Connection struct {
	SessionID  string
	connection net.Conn
}

func CreateServer(port int, u shared.Universe) (s Server) {
	s = Server{status: 0, Playerlist: nil, port: port, u: u}
	return s
}

func (s *Server) AddUser(sessionID string) (result bool) {

	_, prs := s.Playerlist[sessionID]
	if prs == true {
		return false
	}
	s.Playerlist[sessionID] = shared.NewPlayer()
	return true

}

func ConvertMapToBytes(u shared.Universe) ([]byte, error) {
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
