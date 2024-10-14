package main

import (
    "fmt"
    "net"
    "github.com/gdamore/tcell/v2" 
    "log"
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


func main() {
    // Connect to the server
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println(err)
        return
    }

    defer conn.Close()

    _, err = conn.Write([]byte("gommo"))
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
    sessionID := string(buf[:n])
    fmt.Printf("Received SessionID: %s\n", sessionID)

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
        // MARK: Loop
        s.Show()
        ev := s.PollEvent()
        drawBox(s, 0, 0, 50, 50, defStyle, "hello world")

        s.SetContent(0,0, 'H', nil, defStyle)
        s.SetContent(1,0, 'i', nil, defStyle)
        s.SetContent(2,0, 'H', nil, defStyle)

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