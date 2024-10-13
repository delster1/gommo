package main

import (
    "fmt"
	"bytes"
	"encoding/hex"
	"crypto/rand"
    "net"
	"time"
)
type Cell int 

const (
	StateEmpty Cell = iota // makes stateEmpty = 0 , player = 1, statefood = 3, 
	StatePlayer
	StateFood
)
type universe struct {
	Map []Cell
	Width int
	Height int
}
func create_universe(width int, height int) (u universe) {
	u = universe{Map: nil, Width : width, Height : height }
	u.Map = make([]Cell, width * height)
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

	fmt.Println(server_u)
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
		
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte , 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the incoming data
	fmt.Printf("Received: %s\n", buf[:n])
	if !bytes.Equal(buf[:n], []byte("gommo")) {
		fmt.Println("Invalid request")
		conn.Close()
	}
	sessionID, _ := GenerateSessionID(16)
	conn.Write([]byte(sessionID))
	for {
		
		
		response := process_request(buf)	
		if response != nil {
			_, err := conn.Write([]byte(response))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func process_request(buf []byte ) (response []byte) {
	switch {
	case bytes.Equal(buf, []byte("gommo")):
		fmt.Println("Incoming client, generating session id")
	}
	return nil
}