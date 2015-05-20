package main

import (
	"github.com/Lavos/liveproxy"
	"log"
	"flag"
	"bufio"
	"os"
	"fmt"
)

var (
	addr = flag.String("addr", ":11000", "Local address to bind to.")
	lp *liveproxy.LiveProxy
)

func listen() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		err := lp.SwitchTo(scanner.Text())

		if err != nil {
			log.Printf("%s", err)
		} else {
			fmt.Printf("OK\n")
		}
	}
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
