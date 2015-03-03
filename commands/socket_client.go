package main

import (
	"log"
	"net"
)

func main () {
	conn, err := net.Dial("unix", "/tmp/liveproxy.sock")
	defer conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	conn.Write([]byte(":8008\n:8007\nasdfasdf\n"))
	b := make([]byte, 1000)
	n, err := conn.Read(b)
	log.Printf("[%d] `%s` %s", n, b, err)
}
