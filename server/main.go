package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"gommo/shared"
	"net"
	"strconv"
	"strings"
)

func create_universe(size int) (u shared.Universe) {
	u = shared.Universe{Map: nil, Size: size}
	u.Map = make([]shared.Cell, size*size)
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
func iterate_over_cells(u shared.Universe) (new_u shared.Universe) {
	for i := range u.Map {
		// y := i % width
		// x := i % height

		new_u.Map[i] = 0

	}
	return new_u
}
func main() {
	server_u := create_universe(10)
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

// sets up initial connectioConnectionsn
func handleConnection(conn net.Conn, svr *Server) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the incoming data
	fmt.Printf("Received connection: %s\n", buf[:n])

	packetType, err := process_request(buf[:n])
	if err != nil {
		fmt.Printf("Error processing connection request %s\n", err)
		return
	}
	// recieves connection request and sends back sessionid
	err = handle_connection_request(packetType, svr, conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	// recieve map packet
	fmt.Printf("Recieved map: %s\n", buf[:n])
	packettype, err := process_request(buf[:n])
	if err != nil {
		fmt.Println(err)
		return
	}
	// send map packet
	err = handle_connection_request(packettype, svr, conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	// recieve empty map packet
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	packettype, err = process_request(buf[:n])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Recieved final config packet %s\n", buf[:n])
	err = handle_connection_request(packettype, svr, conn)

	for {

		// MARK: Conn Loop
		n, err = conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		// recieved packet from client
		fmt.Printf("Received after validation!: %s\n", buf[:n])
		packettype, err := process_request(buf[:n])
		if err != nil {
			fmt.Println(err)
			return
		}
		// getting response depending on packet
		response, err := handle_request_behavior(packettype, buf, svr)
		if err != nil {
			fmt.Println(err)
			return
		}
		// sending response to client
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Sent response: %s\nendresponse", response)

	}

}

func handle_connection_request(packettype shared.PacketType, svr *Server, con net.Conn) error {
	fmt.Println("Handling connection request")
	switch packettype {
	case shared.PacketTypeConnect:
		sessionid, _ := GenerateSessionID(16)

		svr.AddUser(sessionid)

		base := fmt.Sprintf("gommo\nC\n%s\n", sessionid)
		length := len(base)
		response := fmt.Sprintf("%d\n%s\n", length, base)

		_, err := con.Write([]byte(response))
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	case shared.PacketTypeMove:
		idx := svr.u.Size / 2
		fmt.Println("Recieved empty move packet")
		svr.u.Map[svr.u.Size/2] = shared.User
		base := fmt.Sprintf("gommo\n%c\n%c\n", shared.PacketTypeSuccess, idx)
		packet := fmt.Sprintf("%d\n%s", len(base), base)

		_, err := con.Write([]byte(packet))
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	case shared.PacketTypeMap:
		mapBytes, err := shared.ConvertMapToBytes(svr.u)
		if err != nil {
			fmt.Println(err)
			return err
		}
		finalMapBytes, err := shared.CompressMapData(mapBytes)
		if err != nil {
			fmt.Println(err)
			return err
		}

		base := fmt.Sprintf("gommo\n%c\n%d\n", packettype, svr.u.Size)
		length := len(base)
		packet := fmt.Sprintf("%d\n%s\n", length+len(finalMapBytes), base)
		packet_bytes := []byte(packet)
		map_packet := append(packet_bytes, finalMapBytes...)
		padding := []byte("00000")
		final_packet := append(map_packet, padding...)
		_, err = con.Write(final_packet)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	default:

		errorString := fmt.Sprintf("Incorrect Packet type at connection time %c\n", packettype)
		return errors.New(errorString)
	}
}

func handle_request_behavior(packettype shared.PacketType, buf []byte, svr *Server) ([]byte, error) {
	fmt.Printf("Handling request %c\n", packettype)
	request_str := string(buf)
	parts := strings.Split(request_str, "\n")
	switch packettype {
	case shared.PacketTypeConnect:
		errorString := fmt.Sprintf("Recieved connection packet at incorect time")
		return nil, errors.New(errorString)
	case shared.PacketTypeMap:
		mapBytes, err := shared.ConvertMapToBytes(svr.u)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		finalMapBytes, err := shared.CompressMapData(mapBytes)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		base := fmt.Sprintf("gommo\n%c\n%d\n", packettype, svr.u.Size)
		length := len(base)
		packet := fmt.Sprintf("%d\n%s\n", length+len(finalMapBytes), base)
		packet_bytes := []byte(packet)
		fmt.Printf("FINAL MAP BUTES%b\n", finalMapBytes)
		packet_next := append(packet_bytes, finalMapBytes...)
		final_packet := append(packet_next, []byte("00000000")...)

		// final_packet = append(final_packet, []byte("\n")...)
		return final_packet, nil
	case shared.PacketTypeMove:
		var final_packet []byte

		player := svr.Playerlist[parts[3]]
		playerLocation := player.YPosition*svr.u.Size + player.XPosititon

		svr.u.Map[playerLocation] = shared.Empty
		// user is moving
		idx, err := strconv.Atoi(parts[4])
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		svr.u.Map[idx] = shared.User
		player.XPosititon = idx % svr.u.Size 
		player.YPosition = idx / svr.u.Size

		mapBytes, err := shared.ConvertMapToBytes(svr.u)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		finalMapBytes, err := shared.CompressMapData(mapBytes)
		fmt.Println("FINAL MAP BUTES", finalMapBytes)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		base := fmt.Sprintf("gommo\n%c\n%d", shared.PacketTypeMap, svr.u.Size)
		length := len(base)
		packet := fmt.Sprintf("%d\n%s\n", length+len(finalMapBytes), base)
		packet_bytes := []byte(packet)
		packet_next := append(packet_bytes, finalMapBytes...)
		final_packet = append(packet_next, []byte("00000000")...)

		fmt.Println("FINAL PACKE", final_packet)
		return final_packet, nil
	case shared.PacketTypeDisconnect:
		sessionId := parts[3]
		delete(svr.Playerlist, sessionId)
		return []byte("X"), nil
	default:
		return nil, errors.New("To Be Implemented")
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
