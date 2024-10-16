package main

import (
	"errors"
	"fmt"
	"gommo/shared"
	"log"
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func build_packet(packetType shared.PacketType) []byte {
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
	fmt.Println(response_str)
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

	_, err = conn.Write(build_packet(shared.PacketTypeConnect))
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	sessionID, err := handle_connection_response(buf[:n])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Received SessionID: %s\n", sessionID)

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
	for {
		// drawScreen(s, defStyle)
		// MARK: Loop
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return

		}
		fmt.Printf("Received: %s\n", buf[:n])

		response_str := sessionID + "\nresponse"
		fmt.Println(response_str)

	}

	// Send some data to the server

	// Close the connection
	conn.Close()
}
