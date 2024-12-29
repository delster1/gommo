package main

import (
	"errors"
	"fmt"
	"gommo/shared"
	"log"
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

func initScreen() (tcell.Screen, tcell.Style) {

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault
	s.EnablePaste()
	s.SetStyle(defStyle)
	return s, defStyle
}

func Render(client Client, sigChan chan os.Signal) error {
	s, defStyle := initScreen()
	quit := func() {
		maybePanic := recover()
		s.Fini()
		panic(maybePanic)
	}
	defer quit()

	fmt.Println("INSIDE RENDERER")

	select {
	case <-sigChan:
		return nil
	default:
		update_map(client, s, defStyle) // updates screen with map
		s.Show()
		for {
			s.Show()
			buf := make([]byte, 1024)

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

				} else if ev.Key() == tcell.KeyRight {
					step(&client, shared.Right)

				}
				// fmt.Printf("(%d, %d)\n", client.location.x, client.location.y)
				n, err := client.conn.Write(build_request_packet(client, shared.PacketTypeMove))

				if err != nil {
					fmt.Println(err)
					return err
				}
				_, err = client.conn.Read(buf)
				if err != nil {
					fmt.Println(err)

					return err
				}
				handle_response_behavior(buf[:], n, &client)

			}
			s.Clear()
			update_map(client, s, defStyle)

		}
	}
}

func update_map(client Client, s tcell.Screen, defStyle tcell.Style) {
    for index, cell := range client.u.Map {
        yPos := index / client.u.Size
        xPos := index % client.u.Size

        var renderedCell rune
        switch cell {
        case shared.Empty:
            renderedCell = ' '
        case shared.Land:
            renderedCell = 'L'
        case shared.Water:
            renderedCell = 'W'
        case shared.Mountains:
            renderedCell = 'M'
        case shared.User:
            renderedCell = 'X' // Player is 'X'
        default:
            renderedCell = '?'
        }

        // Render the cell
        s.SetContent(xPos, yPos, renderedCell, nil, defStyle)
    }

    // Ensure the player’s position is rendered explicitly
    playerX := client.location.x
    playerY := client.location.y
    s.SetContent(playerX, playerY, 'X', nil, defStyle)
}

