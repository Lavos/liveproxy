package main

import (
	"log"
	"net"
	"bufio"
)

func main () {
	conn, err := net.Dial("unix", "/tmp/liveproxy.sock")
	defer conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(conn)
	conn.Write([]byte(":8008\n:8007\nasdfasdf\n"))

	var counter int
	for counter < 3 && scanner.Scan() {
		log.Printf("Got Response: %s", scanner.Text())
		counter++
	}

	conn.Close()
}
