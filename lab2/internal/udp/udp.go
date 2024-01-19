package udp

import (
	"fmt"
	"net"
)

func Write(conn net.Conn, m []byte) error{
	_, err := conn.Write(m)
	if err != nil {
		return fmt.Errorf("writing message: %v", err)
	}

	return nil
}

func Read(conn net.Conn) ([]byte, error) {
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		return []byte{}, fmt.Errorf("reading message: %v", err)
	}
	return data[:n], nil
}

func WriteTo(conn net.PacketConn, addr net.Addr, m []byte) error{
	_, err := conn.WriteTo(m, addr)
	if err != nil {
		return fmt.Errorf("writing message: %v", err)
	}
	return nil
}

func ReadFrom(conn net.PacketConn) ([]byte, net.Addr, error) {
	data := make([]byte, 1024)
	n, addr, err := conn.ReadFrom(data)
	if err != nil {
		return []byte{}, nil, fmt.Errorf("reading message: %v", err)
	}
	return data[:n], addr, nil
}
