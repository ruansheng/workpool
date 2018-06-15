package main

import (
	"net"
	"log"
	"workpool"
	"fmt"
)

func main() {
	addr := "127.0.0.1:6789"
	ls, err := net.Listen("tcp", addr)
	if nil != err {
		log.Fatal(err.Error())
	}

	pool, err := workpool.NewPool(10)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		clientConn, err := ls.Accept()
		if err != nil {
			log.Fatal(err.Error())
			continue
		}
		pool.Submit(func() error {
			HandleConnection(&clientConn)
			return nil
		})
	}
}

func HandleConnection(conn *net.Conn) {
	defer func() {
		(*conn).Close()
	}()

	for {
		buf := make([]byte, 1024)
		n, err := (*conn).Read(buf)
		if err != nil {
			break
		}
		fmt.Println(string(buf[:n]))

		result := "hello world"
		(*conn).Write([]byte(result))
	}
}