package main

import (
	"errors"
	"github.com/gdamore/tcell/v2"
	"gommo/shared"
	"log"
	"os"
	"os/signal"
	"syscall"
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

func DrawScreen(s tcell.Screen, defStyle tcell.Style, size int) {
	s.Show()
	s.Sync()
}

func NewScreen() (tcell.Screen, tcell.Style) {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	// s.EnableMouse()
	s.EnablePaste()
	// Clear screen
	s.Clear()
	return s, defStyle
}

func Render(client Client, defStyle tcell.Style, s tcell.Screen, sigChan chan os.Signal) error {
	size := client.u.Size

	cMap := client.u.Map
	var newMap shared.Universe
	newMap.Size = size
	DrawScreen(s, defStyle, size) // Draw the initial screen
	//	drawBox(s, 0, 0, size, size, defStyle, "")
	drawBox(s, 0, 0, client.u.Size, client.u.Size, defStyle, "")

	select {
	case <-sigChan:
		return nil
	default:
		for index, cell := range cMap {
			if (index == 0 || index % client.u.Size == 0){
				continue
			}
			s.Show()

			yPos := index / size
			xPos := index % size

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
		s.Show()
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
					return errors.New("Exiting Screen")
				}
			}
		}
	}
	return nil
}
