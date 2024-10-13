package main

import (
    "fmt"
    "net"
)
type Cell int 

const (
	StateEmpty Cell = iota
	StatePlayer
	StateFood
)
type universe struct {
	Map []Cell
}

func create_universe() (u universe) {
	u = universe{Map: nil}
	u.Map = make([]Cell, 100)
	return u
} 

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	for {
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
	_, err := conn.Read(buf)
	if err != nil {
        fmt.Println(err)
        return
    }

    // Print the incoming data
    fmt.Printf("Received: %s", buf)
}