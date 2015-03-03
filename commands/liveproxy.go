package main

import (
	"github.com/Lavos/liveproxy"
	"log"
	"time"
)

func main(){
	lp, err := liveproxy.New(":11000")

	log.Printf("Liveproxy: %#v, %#v", lp, err)

	lp.SwitchTo(":8007")

	go func(){
		var toggle bool
		t := time.NewTicker(10 * time.Second)

		for {
			<-t.C
			toggle = !toggle

			if toggle {
				lp.SwitchTo(":8008")
			} else {
				lp.SwitchTo(":8007")
			}
		}
	}()

	<-lp.Done
}
