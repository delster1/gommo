package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"gommo/shared"
	"net"
	"strings"
)

type Cell int

type universe struct {
	Map    []Cell
	Width  int
	Height int
}

func create_universe(width int, height int) (u universe) {
	u = universe{Map: nil, Width: width, Height: height}
	u.Map = make([]Cell, width*height)
	return u
}

func GenerateSessionID(length int) (string, error) {
	// Create a byte slice with the specified length
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err // Return an error if random generation fails
	}

	// Return the session ID as a hex string
	return hex.EncodeToString(bytes), nil
}
func iterate_over_cells(u universe) (new_u universe) {
	for i := range u.Map {
		// y := i % width
		// x := i % height

		new_u.Map[i] = 0

	}
	return new_u
}

func main() {
	server_u := create_universe(50, 50)
	svr := CreateServer(8080, server_u)

	fmt.Println(svr.u.Map)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	for {
		//MARK: server loop
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn, &svr)
	}
}

func handleConnection(conn net.Conn, svr *Server) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the incoming data
	fmt.Printf("Received: %s\n", buf[:n])

	packetType, err := process_request(buf[:n])
	if err != nil {
		fmt.Printf("Error processing connection request %s\n", err)
		return
	}
	handle_connection_request(packetType, svr, conn)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Recieved: %s\n", buf[:n])
		packettype, err := process_request(buf[:n])

		handle_request_behavior(packettype, buf, svr)
	}
}

func handle_connection_request(packettype shared.PacketType, svr *Server, con net.Conn) error {
	switch packettype {
	case shared.PacketTypeConnect:
		sessionid, _ := GenerateSessionID(16)
		connection := Connection{sessionid, con}
		svr.Connections = append(svr.Connections, connection)
		base := fmt.Sprintf("gommo\nC\n%s\n", sessionid)
		length := len(base)
		response := fmt.Sprintf("%d\n%s\n", length, base)
		con.Write([]byte(response))
		return nil
	default:
		errorString := fmt.Sprintf("Incorrect Packet type at connection time %c\n", packettype)
		return errors.New(errorString)
	}
}
func handle_request_behavior(packettype shared.PacketType, buf []byte, svr *Server) error {
	switch packettype {
	case shared.PacketTypeConnect:
		errorString := fmt.Sprintf("Recieved connection packet at incorect time")
		return errors.New(errorString)
	default:
		return errors.New("To Be Implemented")
	}
}

func process_request(request []byte) (shared.PacketType, error) {
	request_str := string(request)
	parts := strings.Split(request_str, "\n")
	errorString := fmt.Sprintf("Invalid Packet:\t%s", string(request))
	if parts[1] != "gommo" {
		return shared.PacketTypeErr, errors.New(errorString)
	}

	switch parts[2] {
	case "C":
		return shared.PacketTypeConnect, nil
	case "M":
		return shared.PacketTypeMap, nil
	case "L":
		return shared.PacketTypeMove, nil
	case "X":
		return shared.PacketTypeDisconnect, nil
	case "E":
		return shared.PacketTypeErr, nil
	default:
		return shared.PacketTypeErr, nil
	}
}
