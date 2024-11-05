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

func build_request_packet(packetType shared.PacketType) []byte {
	switch packetType {
	case shared.PacketTypeConnect:
		base := fmt.Sprintf("gommo\n%c\n", packetType)
		packetLen := len(base)
		packet := fmt.Sprintf("%d\n%s", packetLen, base)
		return []byte(packet)
	case shared.PacketTypeMap:
		// MARK: TBI
		base := fmt.Sprintf("gommo\n%c\n", packetType)
		packetLen := len(base)
		packet := fmt.Sprintf("%d\n%s", packetLen, base)
		return []byte(packet)
	case shared.PacketTypeMove:
		// TODO:
		return []byte("gommo\nE\n")
	case shared.PacketTypeDisconnect:
		return []byte("X")
	default:
		return []byte("gommo\nE\n")
	}
}

func handle_connection_response(response []byte) (string, error) {
	response_str := string(response)
	packet_parts := strings.Split(response_str, "\n")
	if packet_parts[1] != "gommo" {
		fmt.Println("bad packet recieved")
		return "", errors.New("Bad Packet Recieved to Connection")
	}
	switch packet_parts[2] {
	case "C":
		sessionID := packet_parts[3]
		return sessionID, nil
	default:
		errorString := fmt.Sprintf("Recieved incorrect packet to handle connection's response, %s\n", packet_parts[2])
		return "", errors.New(errorString)
	}
}
func handle_response_behavior(response []byte) (interface{}, error) {
	response_str := string(response)
	packet_parts := strings.Split(response_str, "\n")
	// packet_len := packet_parts[0]
	if packet_parts[1] != "gommo" || packet_parts[2] == "C" {
		errorString := fmt.Sprintf("bad packet recieved %b", packet_parts[1])
		return "", errors.New(errorString)
	}
	switch packet_parts[2] {
	case "M":
		fmt.Printf("Received Map Packet: %b\n", response)
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
	s, defStyle := NewScreen()
	var screenPointer *tcell.Screen = &s
	defer conn.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	client := BuildClient(conn, screenPointer, errChan) // builds client w/ sessionid & stuff given connection
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
			n, err := conn.Read(buf)
			if err != nil {
				errChan <- err
				return
			}
			fmt.Printf("Received from server:%b\n", buf[:n])
		}
	}()
	fmt.Println("running renderer")
	err = Render(client, defStyle, *client.screen, sigChan)
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
