package main

import (
	"net"
	"log"
	"code.qschou.com/dbj-hd/workpool"
	"strings"
	"strconv"
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

		lines := ParseProtocal(string(buf[:n]))

		if len(lines) == 4 && strings.EqualFold(lines[2], "COMMAND") {
			pong := "+OK\r\n"
			(*conn).Write([]byte(pong))
			continue
		}

		if len(lines) == 4 && strings.EqualFold(lines[2], "PING") {
			pong := "+PONG\r\n"
			(*conn).Write([]byte(pong))
			continue
		}

		if len(lines) !=6 {
			pong := "+不支持的命令\r\n"
			(*conn).Write([]byte(pong))
			continue
		}

		if !strings.EqualFold(lines[2], "GET") {
			pong := "-不支持的命令\r\n"
			(*conn).Write([]byte(pong))
			continue
		}

		result := "hello world"
		redisCmd := MakeGetProtocal(result)
		(*conn).Write([]byte(redisCmd))
	}
}

func ParseProtocal(cmd string) []string {
	lines := strings.Split(cmd, "\r\n")
	return lines
}

func MakeGetProtocal(resp string) string {
	llen := len(resp)
	lines := make([]string, 5, 5)
	lines = append(lines, "$")
	lines = append(lines, strconv.Itoa(llen))
	lines = append(lines, "\r\n")
	lines = append(lines, resp)
	lines = append(lines, "\r\n")
	return strings.Join(lines, "")
}
