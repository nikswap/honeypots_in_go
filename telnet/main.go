package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	addr := ":2323"
	server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer server.Close()

	log.Println("Server is running on:", addr)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Failed to accept conn.", err)
			continue
		}

		go func(conn net.Conn) {
			defer func() {
				conn.Close()
			}()
			n, err := conn.Write([]byte("$"))
			t := time.Now()
			CheckError(err)
			fmt.Println("Bytes written", n)
			lines_sent := make([]string, 100)
			buf := make([]byte, 0, 4096)
			tmp := make([]byte, 256)
			lb_counter := 0
			for {
				n, err := conn.Read(tmp)
				if err != nil {
					if err != io.EOF {
						fmt.Println("read error:", err)
					}
					break
				}
				buf = append(buf, tmp[:n]...)
				old_lb := lb_counter
				for i, c := range buf[lb_counter:] {
					if c == 0xA {
						lb_counter = i
						break
					}
				}
				if old_lb < lb_counter {
					lines_sent = append(lines_sent, base64.StdEncoding.EncodeToString(buf[old_lb:lb_counter]))
					fmt.Println(t.Format("2006-01-02 15:04:05"), ",", conn.RemoteAddr().String(), ",", lines_sent[len(lines_sent)-1])
					conn.Write([]byte("\n$"))
				}
			}
			fmt.Println("Connection done")
		}(conn)
	}
}
