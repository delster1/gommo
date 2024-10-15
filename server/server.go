package main

import "gommo/shared"

type server struct {
	status     int
	Playerlist map[string]*shared.Player
	port       int
	u          universe
}

func create_server(port int, u universe) (s server) {
	s = server{status: 0, Playerlist: nil, port: port, u: u}
	return s
}

func (s *server) AddUser(sessionID string) (result bool) {

	_, prs := s.Playerlist[sessionID]
	if prs == true {
		return false
	}
	s.Playerlist[sessionID] = shared.NewPlayer()
	return true

}
