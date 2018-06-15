package main

import (
	"net"
	"fmt"
	"io"
	"os"
	"log"
)

func main()  {
	addr := "127.0.0.1:6789"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	for i:=0; i < 100; i++ {
		conn.Write([]byte("test"))
	}

	go Handle(conn)

	for{
		select {}
	}
}

func Handle(conn net.Conn) {
	for {
		data := make([]byte, 1024)
		buf := make([]byte, 128)
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
			}
			data = append(data, buf[:n]...)
			if n != 128 {
				break
			}
		}

		fmt.Println(string(data))
	}
}