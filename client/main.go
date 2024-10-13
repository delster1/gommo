package main

import (
    "fmt"
    "net"
)

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

    for {
        // MARK: client Loop
        // this is where operations from the client are sent and where we will draw the window and stuff
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