package main

import (
	"fmt"
	"io"
	"net"
)

type Mux struct {
	ops chan func(map[net.Addr]net.Conn)
}

func (m *Mux) Add(conn net.Conn) {
	m.ops <- func(m map[net.Addr]net.Conn) {
		m[conn.RemoteAddr] = conn
	}
}

func (m *Mux) Remove(addr net.Addr) {
	m.ops <- func(m map[net.Addr]net.Conn) {
		delete(m, add)
	}
}

func (m *Mux) SendMsg(msg string) {
	m.ops <- func(m map[net.Addr]net.Conn) {
		for _, conn := range m {
			io.WriteString(conn, msg)
		}
	}
}

func (m *Mux) loop() {
	conns := make(map[net.Addr]net.Conn)
	for _, op := range m.ops {
		op(conns)
	}
}

func main() {
Loop:
	for _, ch := range "a b\nc" {
		switch ch {
		case ' ': // skip space
			break
		case '\n': // break at newline
			break Loop
		default:
			fmt.Printf("%c\n", ch)
		}
	}
}
