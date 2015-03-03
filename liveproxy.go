package liveproxy

import (
	"io"
	"log"
	"net"
)

type LiveProxy struct {
	Listener *net.TCPListener

	Done           chan struct{}
	Control        chan *net.TCPAddr
	GetDestination chan *net.TCPAddr
}

func New(local string) (*LiveProxy, error) {
	local_address, err := net.ResolveTCPAddr("tcp4", local)

	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp4", local_address)

	if err != nil {
		return nil, err
	}

	l := &LiveProxy{
		Listener:       listener,
		Control:        make(chan *net.TCPAddr),
		GetDestination: make(chan *net.TCPAddr),
		Done:           make(chan struct{}),
	}

	go l.control()
	go l.connect()

	log.Printf("Running on: %s", local_address)

	return l, nil
}

func (l *LiveProxy) control() {
	var destination *net.TCPAddr

	for {
		select {
		case destination = <-l.Control:
		case l.GetDestination <- destination:
		}
	}
}

func (l *LiveProxy) connect() {
	for {
		var conn *net.TCPConn
		var dest *net.TCPConn
		var destination_address *net.TCPAddr
		var err error

		log.Printf("Waiting for connection...")
		conn, err = l.Listener.AcceptTCP()
		log.Printf("Connection attempted.")

		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}

		destination_address = <-l.GetDestination

		dest, err = net.DialTCP("tcp4", nil, destination_address)

		if err != nil {
			log.Printf("Could not connect to Destination: %s", destination_address)
			conn.Close()
			continue
		}

		a := make(chan struct{})
		b := make(chan struct{})

		go l.pipe(conn, dest, a)
		go l.pipe(dest, conn, b)

		go func() {
			select {
			case <-a:
				log.Printf("got close from a")
				conn.CloseRead()

			case <-b:
				log.Printf("got close from b")
				dest.CloseRead()
			}
		}()
	}
}

func (l *LiveProxy) pipe(dest, src *net.TCPConn, signal chan struct{}) {
	io.Copy(dest, src)
	src.Close()
	signal <- struct{}{}
}

func (l *LiveProxy) SwitchTo(destination_address string) error {
	address, err := net.ResolveTCPAddr("tcp4", destination_address)

	if err != nil {
		return err
	}

	log.Printf("Destination address is now: %s", address)
	l.Control <- address

	return nil
}
