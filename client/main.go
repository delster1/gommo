package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"gommo/shared"
	"io"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gdamore/tcell/v2"
)

func DecompressMapData(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	reader, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	_, err = io.Copy(&out, reader)
	if err != nil {
		return nil, err
	}
	reader.Close()
	return out.Bytes(), nil
}
func ConvertBytesToMap(data []byte) (shared.Universe, error) {

	temp := int(len(data))
	u := shared.Universe{
		Map:    make([]shared.Cell, temp),
		Width:  int(math.Sqrt(float64(temp))),
		Height: int(math.Sqrt(float64(temp))),
	}
	fmt.Println("created universe")
	buf := bytes.NewReader(data)
	for i := range u.Map {
		if err := binary.Read(buf, binary.LittleEndian, &u.Map[i]); err != nil {
			return u, err
		}
	}
	return u, nil
}

func build_request_packet(packetType shared.PacketType) []byte {
	switch packetType {
	case shared.PacketTypeConnect:
		base := fmt.Sprintf("gommo\n%c\n", packetType)
		packetLen := len(base)
		packet := fmt.Sprintf("%d\n%s", packetLen, base)
		return []byte(packet)
	case shared.PacketTypeMap:

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

func handle_response_behavior(response []byte) (string, error) {
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
		mapDataStart := strings.Index(response_str, "\nM\n") + 3
		universeBytes, err := DecompressMapData(response[mapDataStart:])
		if err != nil {
			errStr := fmt.Sprintf("Error decompressing map, %s", err)
			fmt.Println(errStr)
			return "", err
		}
		fmt.Printf("decompressed Result%b\n", universeBytes)
		clientUniverse, err := ConvertBytesToMap(universeBytes)
		if err != nil {
			errStr := fmt.Sprintf("Error converting bytes to map %s", err)
			fmt.Println(errStr)
			return "", err
		}
		fmt.Println(clientUniverse.Map)
		return string(clientUniverse.Map), nil
	default:
		return "to be implemented", errors.New("To Be Implemented")
	}

}
func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}
func drawScreen(s tcell.Screen, defStyle tcell.Style) {
	s.Show()
	ev := s.PollEvent()
	drawBox(s, 0, 0, 50, 50, defStyle, "hello world")

	s.SetContent(0, 0, 'H', nil, defStyle)
	s.SetContent(1, 0, 'i', nil, defStyle)
	s.SetContent(2, 0, 'H', nil, defStyle)

	switch ev := ev.(type) {
	case *tcell.EventResize:
		s.Sync()
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
			return
		} else if ev.Key() == tcell.KeyCtrlL {
			s.Sync()
		} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
			s.Clear()
		}
	}

}
func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		buf := make([]byte, 1024)
		_, err = conn.Write(build_request_packet(shared.PacketTypeConnect))
		if err != nil {
			fmt.Println(err)
			return
		} // sends connection request

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		} // loads buf with connection info from server

		_, err = handle_connection_response(buf[:n])
		if err != nil {
			fmt.Println(err)
			return
		} // grabs sessionid from buf

		//MARK: Testing stuff here rn
		_, err = conn.Write(build_request_packet(shared.PacketTypeMap))
		if err != nil {
			fmt.Println(err)
			return
		}
		n, err = conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		server_map, err := handle_response_behavior(buf[:n])
		fmt.Println(server_map)
		for {
			//MARK: Main looping handler
			n, err := conn.Read(buf)
			if err != nil {
				errChan <- err
				return
			}
			fmt.Printf("Received from server:%b\n", buf[:n])
		}
	}()
	for {
		select {
		case <-sigChan:
			// Gracefully handle the signal (Ctrl+C or termination)
			fmt.Println("Exiting, closing connection...")
			conn.Close()
			return
		case err := <-errChan:
			// Handle any network errors
			fmt.Println("Network error:", err)
			return
		}
	}

	// crete new screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()

	s.Clear()
	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()
	// for {
	// 	// drawScreen(s, defStyle)
	// 	// MARK: Loop
	// 	n, err := conn.Read(buf)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return

	// 	}
	// 	fmt.Printf("Received: %s\n", buf[:n])

	// 	response_str := sessionID + "\nresponse"
	// 	fmt.Println(response_str)

	// }

	// Send some data to the server

	// Close the connection
	conn.Close()
}
