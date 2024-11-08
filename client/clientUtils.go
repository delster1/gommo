package main

import (
	"fmt"
	"gommo/shared"
	"net"
	"strings"
	"strconv"

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
}
func step(c *Client, direction shared.Dir) {
    old_idx := c.u.Size * c.location.y + c.location.x

    // Update the location based on direction
    switch direction {
    case shared.Up:
        c.location.y = (c.location.y - 1 + c.u.Size) % c.u.Size
    case shared.Down:
        c.location.y = (c.location.y + 1) % c.u.Size
    case shared.Left:
        c.location.x = (c.location.x - 1 + c.u.Size) % c.u.Size
    case shared.Right:
        c.location.x = (c.location.x + 1) % c.u.Size
    default:
        return // No movement
    }

    // Calculate the new index after movement
    idx := c.u.Size * c.location.y + c.location.x

    // Update the map logically
    c.u.Map[old_idx] = shared.Empty   // Set the old cell to Empty
    c.u.Map[idx] = shared.User        // Set the new cell to User
}
func handle_response_behavior(response []byte, c *Client) {
	response_str := string(response)
	packet_parts := strings.Split(response_str, "\n")

	switch packet_parts[1] {
	case "M":
		var serverUniverse shared.Universe

		mapDataStart := strings.Index(response_str, "\nM\n") + 7
		universeBytes, err := shared.DecompressMapData(response[mapDataStart:])
		if err != nil {
			errStr := fmt.Sprintf("Error decompressing map, %s", err)
			fmt.Println(errStr)
			return 
		}
		mapSize, err := strconv.Atoi(packet_parts[3])
		if err != nil {
			fmt.Println(err)
			return
		}
		serverUniverse, err = shared.ConvertBytesToMap(mapSize, universeBytes)
		if err != nil {
			errStr := fmt.Sprintf("Error converting bytes to map %s", err)
			fmt.Println(errStr)
			return
		}
		// Update the client's universe with the server's universe 
		for i, cell := range serverUniverse.Map {
			if c.u.Map[i] != cell {

				c.u.Map[i] = cell
			}
		}

			
		
	case "L":
		panic("Recieved move packet from server???")	
	default:
		fmt.Println("Unknown packet type")
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



func BuildClient(conn net.Conn,errChan chan<- error) (c Client) {
	var universe shared.Universe
	buf := make([]byte, 1024)
	_, err := conn.Write(build_configure_packet(shared.PacketTypeConnect))
	if err != nil {
		fmt.Println(err)
		return
	} // sends connection request

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	} // loads buf with connection info from server

	result, err := handle_setup_behavior(buf[:n])
	if err != nil {
		fmt.Println(err)
		return
	} // grabs sessionid from buf

	sessionID, ok := result.(string) 
	if !ok { 
		fmt.Println("Error converting sessionID")
		return
	}
	_, err = conn.Write(build_configure_packet(shared.PacketTypeMap))
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
	fmt.Println("GOT MAP")

	u, err := handle_setup_behavior(buf[:n])
	// handle map response - convert to universe
	if univ, ok := u.(shared.Universe); ok {
		universe = univ
	} else {
		fmt.Println(err)
		return
	}
	

	// send move packet with client's initial location (always midpoint)	
	_, err = conn.Write(build_configure_packet(shared.PacketTypeMove))
	if err != nil {
		errString := fmt.Sprintf("Error handling map %s\n", err)
		fmt.Println(errString)
		errChan <- err
		return
	}
	fmt.Println("SENT MOVE PACKET")
	// recieve and check move packet validation
	print("BUILT CLIENT")
	c = Client{u: shared.Universe{Map: universe.Map, Size: universe.Size}, location: Location{x: universe.Size / 2 , y: universe.Size / 2}, SessionID: string(sessionID), conn: conn, isConnected: true }
	idx := c.u.Size * c.location.y + c.location.x
	c.u.Map[idx] = 'P' 
	

	return c
}
