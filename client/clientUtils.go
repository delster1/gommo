package main

import (
	"fmt"
	"gommo/shared"
	"net"

	"github.com/gdamore/tcell/v2"
)

type Location struct {
	x int
	y int
}

type Client struct {
	u           shared.Universe
	location    Location
	SessionID   string
	conn        net.Conn
	isConnected bool
	screen      *tcell.Screen
}
func step(direction Dir){
	switch direction:
	case Up{
	} 
}
func BuildClient(conn net.Conn, screen *tcell.Screen, errChan chan<- error) (c Client) {
	var universe shared.Universe
	buf := make([]byte, 1024)
	_, err := conn.Write(build_request_packet(shared.PacketTypeConnect))
	if err != nil {
		fmt.Println(err)
		return
	} // sends connection request

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	} // loads buf with connection info from server

	sessionID, err := handle_connection_response(buf[:n])
	if err != nil {
		fmt.Println(err)
		return
	} // grabs sessionid from buf

	_, err = conn.Write(build_request_packet(shared.PacketTypeMap))
	// sending map over
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err = conn.Read(buf)
	// recieve map response
	if err != nil {
		fmt.Println(err)
		return
	}
	u, err := handle_response_behavior(buf[:n])
	// handle map response - convert to universe
	if univ, ok := u.(shared.Universe); ok {
		universe = univ
	}
	// universe returned from server as universe object

	if err != nil {
		errString := fmt.Sprintf("Error handling map %s\n", err)
		fmt.Println(errString)
		errChan <- err
		return
	}
	c = Client{u: shared.Universe{Map: universe.Map, Size: universe.Size}, location: Location{x: universe.Size / 2 , y: universe.Size / 2}, SessionID: sessionID, conn: conn, isConnected: true, screen: screen}
	idx := c.u.Size * c.location.y + c.location.x
	c.u.Map[idx] = 'P' 
	

	return c
}
