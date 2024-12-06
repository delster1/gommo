package main

import (
	"gommo/shared"
	"net"
)

type Server struct {
	status      bool
	Playerlist  map[string]*shared.Player
	port        int
	u           shared.Universe
}
type Connection struct {
	SessionID  string
	connection net.Conn
}

func CreateServer(port int, u shared.Universe) (s Server) {
	s = Server{status: true, Playerlist: nil, port: port, u: u }
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
