package main

import (
	"errors"
	"log"
	"fmt"
	"gommo/shared"
"os"
	"os/signal"
	"syscall"

	"github.com/gdamore/tcell/v2"
)

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


func NewScreen() (tcell.Screen, tcell.Style, error) {
    s, err := tcell.NewScreen()
    if err != nil {
        return nil, tcell.Style{}, fmt.Errorf("failed to create new screen: %w", err)
    }
    if err := s.Init(); err != nil {
        return nil, tcell.Style{}, fmt.Errorf("failed to initialize screen: %w", err)
    }
    defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
    s.SetStyle(defStyle)
    s.EnablePaste()
    s.Clear()
    return s, defStyle, nil
}




func Render(client Client,  sigChan chan os.Signal) error {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	defStyle := tcell.StyleDefault

	s.SetStyle(defStyle)
	drawBox(s, 0, 0, client.u.Size, client.u.Size, defStyle, "")
	quit := func() {
		maybePanic := recover()
		s.Fini()
			panic(maybePanic)
	}
	defer quit()
		
	size := client.u.Size
	fmt.Println("INSIDE RENDERER")
	var newMap shared.Universe
	newMap.Size = size
	//	drawBox(s, 0, 0, size, size, defStyle, "")
	

	select {
	case <-sigChan:
		return nil
	default:
		update_map(client, s, defStyle) // updates screen with map	
		s.Show()
		//buf := make([]byte, 1024)
		for {
			fmt.Println("INSIDE RENDERER LOOP")
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
					return errors.New("Exiting Screen")
				} else if ev.Key() == tcell.KeyUp {
					step(&client, shared.Up)
				} else if ev.Key() == tcell.KeyDown {
					step(&client, shared.Down)

				} else if ev.Key() == tcell.KeyLeft {
					step(&client, shared.Left)

				} else if ev.Key() == tcell.KeyRight{
					step(&client, shared.Right)
					
				}

				build_request_packet(client, shared.PacketTypeMove)
				//_, err := client.conn.Write(build_request_packet(client, shared.PacketTypeMove))
				//if err != nil {
				//	fmt.Println(err)
				//	return err
				//}
				//n, err := client.conn.Read(buf)
				//if err != nil {
				//	fmt.Println(err)
				//	return err
				
				//handle_response_behavior(buf[:n], &client)
				
			}
			fmt.Println(client.location)
			update_map(client, s, defStyle)
			
		}
	}
	return nil
}
// function to iterate through the map and update the screen given what is at each cell  
func update_map(client Client, s tcell.Screen, defStyle tcell.Style){ 
	for index, cell := range client.u.Map {
		if (index == 0 || index % client.u.Size == 0){
			continue
		}

		yPos := index / client.u.Size
		xPos := index % client.u.Size

		var renderedCell rune

		switch cell {
		case 0:
			renderedCell = ' '
		case 1:
			renderedCell = 'L' // make these fancy runes later
		case 2:
			renderedCell = 'W'
		case 3:
			renderedCell = 'M'
		case 4:
			renderedCell = 'P'
		default:
			renderedCell = 'X'
		}
		
		s.SetContent(xPos, yPos, renderedCell, nil, defStyle)
	}
}
