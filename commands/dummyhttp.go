package main

import (
        "net/http"
        "flag"
        "log"
)

var (
        port = flag.String("port", ":8000", "Listen at address specified.")
)

func noop (w http.ResponseWriter, req *http.Request) {}

func main () {
        flag.Parse()
        log.Printf("Listening at %s.", *port)

        s := &http.Server{
                Addr:           *port,
                Handler:        http.HandlerFunc(noop),
                MaxHeaderBytes: 1 << 20,
        }

        log.Fatal(s.ListenAndServe())
}
