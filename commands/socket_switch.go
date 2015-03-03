package main

import (
	"net"
	"github.com/Lavos/liveproxy"
	"log"
	"flag"
	"bufio"
	"os"
)

var (
	addr = flag.String("addr", ":11000", "Local address to bind to.")
	lp *liveproxy.LiveProxy
)

func listen() {
	socket, err := net.Listen("unix", "/tmp/liveproxy.sock")
	defer os.Remove("/tmp/liveproxy.sock")

	if err != nil {
		log.Fatal(err)
	}

	var conn net.Conn

	for {
		conn, err = socket.Accept()

		if err != nil {
			log.Print(err)
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		log.Printf("GOT `%s`", scanner.Text())
		err := lp.SwitchTo(scanner.Text())

		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
			continue
		}

		conn.Write([]byte("OK\n"))
	}

	log.Printf("Scanning complete.")
	conn.Close()
}

func main() {
	flag.Parse()

	var err error
	lp, err = liveproxy.New(*addr)

	if err != nil {
		log.Fatal(err)
	}

	go listen()
	<-lp.Done
}
