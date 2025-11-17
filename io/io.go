package io

import (
	"net"

	"github.com/dhairyajoshi/gomq/parsers"
)

var parser = parsers.GetParser()

type IOHandler struct {
	conn *net.Conn
}

func (hand *IOHandler) Write(resp parsers.ServerResponse) (n int, err error) {
	n, err = (*hand.conn).Write(parser.Encode(resp))
	return
}

func (hand *IOHandler) Read() ([]byte, error) {
	buff := make([]byte, 1024)
	n, err := (*hand.conn).Read(buff)
	return buff[:n], err
}

func (hand *IOHandler) Close() (err error) {
	err = (*hand.conn).Close()
	return
}

func NewIoHandler(conn *net.Conn) *IOHandler {
	return &IOHandler{conn}
}
