package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Transporter interface {
	ListenAndAccept() error
}

type TCPTransporter struct {
	listenAddr string
	listener net.Listener
}

func NewTCPTransporter(listenAddr string) *TCPTransporter {
	return &TCPTransporter{
		listenAddr: listenAddr,
	}
}

func (t *TCPTransporter) ListenAndAccept() error {
	var err error 
	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptingLoop()

	return nil
}

func (t *TCPTransporter) startAcceptingLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error %s\n", err)
		}

		go t.handleConnection(conn)
	}
}

func (t *TCPTransporter) handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("New incoming connection from %v\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Connection error: %v\n", err)
			break
		}
		fmt.Printf("Received message: %s", message)
	}
}

func main() {
	t := NewTCPTransporter(":3000")
	if err := t.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}