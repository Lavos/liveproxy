package main

import (
	"net"
	"bufio"
	"flag"
	"log"
)

func main () {
	flag.Parse()
	addr := flag.Arg(0)

	if addr == "" {
		log.Fatal("No address found.")
	}

	conn, err := net.Dial("unix", "/tmp/liveproxy.sock")
	defer conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(conn)
	conn.Write([]byte(addr + "\n"))

	ok := scanner.Scan()

	if ok {
		log.Printf("LiveProxy returned: %s", scanner.Text())
	} else {
		log.Printf("No response.")
	}
}
