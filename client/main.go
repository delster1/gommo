package main

import (
	"errors"
	"fmt"
	"gommo/shared"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gdamore/tcell/v2"
)
func build_configure_packet( packetType shared.PacketType) []byte {
	switch packetType {
	case shared.PacketTypeConnect:
		base := fmt.Sprintf("gommo\n%c\n", packetType)
		packetLen := len(base)
		packet := fmt.Sprintf("%d\n%s\n", packetLen, base)
		return []byte(packet)
	case shared.PacketTypeMap:
		// MARK: TBI
		base := fmt.Sprintf("gommo\n%c\n", packetType)
		packetLen := len(base)
		packet := fmt.Sprintf("%d\n%s\n", packetLen, base)
		return []byte(packet)
	case shared.PacketTypeMove:
		// TODO:
		base := fmt.Sprintf("gommo\n%c\n", packetType)
		packet := fmt.Sprintf("%d\n%s\n", len(base), base)
		fmt.Println(len(base))
		return []byte(packet)
	case shared.PacketTypeDisconnect:
		return []byte("X")
	default:
		return []byte("gommo\nE\n")
	}
}
func build_request_packet(c Client, packetType shared.PacketType) []byte {
	switch packetType {
	case shared.PacketTypeConnect:
		panic("Connect packet should not be sent by client at this stage")
	case shared.PacketTypeMap:
		panic("User shouldn't be sending map requests - only move requests!")
	case shared.PacketTypeMove:
		// custom move packet - base, sessionid, x, y
		base := fmt.Sprintf("gommo\n%c\n%s\n%d\n%d\n", packetType, c.SessionID, c.location.x, c.location.y)
		packet := fmt.Sprintf("%d\n%s", len(base), base)
		return []byte(packet)
	default:
		return []byte("gommo\nE\n")
	}
}


func handle_setup_behavior(response []byte) (interface{}, error) {
	response_str := string(response)
	packet_parts := strings.Split(response_str, "\n")
	// packet_len := packet_parts[0]
	if packet_parts[1] != "gommo" {
		errorString := fmt.Sprintf("bad packet recieved %s",response_str)
		return "", errors.New(errorString)
	}
	switch packet_parts[2] {

	case "M": // map packet
		mapDataStart := strings.Index(response_str, "\nM\n") + 7
		universeBytes, err := shared.DecompressMapData(response[mapDataStart:])
		if err != nil {
			errStr := fmt.Sprintf("Error decompressing map, %s", err)
			fmt.Println(errStr)
			return "", err
		}
		mapSize, err := strconv.Atoi(packet_parts[3])
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		clientUniverse, err := shared.ConvertBytesToMap(mapSize, universeBytes)
		if err != nil {
			errStr := fmt.Sprintf("Error converting bytes to map %s", err)
			fmt.Println(errStr)
			return "", err
		}
		return clientUniverse, nil
	case "L": // client recieved move packet
		if packet_parts[0] == "8" {
			return "", nil	
		} else {
			return "", nil
		}
	case "C":
		var sessionID string
		sessionID = packet_parts[3]
		return sessionID, nil

	default:
		return "to be implemented", errors.New("To Be Implemented")
	}

}

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	s, err := tcell.NewScreen()
	defStyle := tcell.StyleDefault
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("HELOOOOO")
	defer conn.Close()


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	errChan := make(chan error, 1)
	// BuildClient collects sessionID, universe, and other info from server
	client := BuildClient(conn, s, errChan)
	fmt.Println(client)
	buf := make([]byte, 1024)

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()
		go func() {
		for {
			x, err := conn.Read(buf)
			if err != nil {
				errChan <- err
				return
			}
			fmt.Printf("QUITTING from server:%b\n", buf[:x])
		}
	}()
	fmt.Println("running renderer")
	//MAIN LOOP IN RENDER
	err = Render(client, defStyle, client.screen, sigChan)
	if err != nil {
		return
	}
	fmt.Println("renderer ended")
	for {
	select {
		case <-sigChan:
			// Gracefully handle the signal (Ctrl+C or termination)
			fmt.Println("Exiting, closing connection...")
			conn.Close()
			return
		case err := <-errChan:
			// Handle any network errors or errors from Render
			fmt.Println("Error:", err)
			return
		}
	}
}
