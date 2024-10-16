package main

import (
	"gommo/shared"
	"net"
)

type Server struct {
	status      int
	Playerlist  map[string]*shared.Player
	port        int
	u           universe
	Connections []Connection
}
type Connection struct {
	SessionID  string
	connection net.Conn
}

func CreateServer(port int, u universe) (s Server) {
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
