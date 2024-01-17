package udp

import (
	"fmt"
	"lab2/internal/message"
	"net"
)

func Write(conn net.Conn, m message.Message) error{
	data, err := message.ToBytes(m)
	if err != nil {
		return fmt.Errorf("converting message to bytes: %v", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("writing message: %v", err)
	}

	return nil
}

func Read(conn net.Conn, m message.Message) error {
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		return fmt.Errorf("reading message: %v", err)
	}

	err = message.Parse(m, data[:n])
	if err != nil {
		return fmt.Errorf("parsing message: %v", err)
	}

	return nil
}

func WriteTo(conn net.PacketConn, addr net.Addr, m message.Message) error{
	data, err := message.ToBytes(m)
	if err != nil {
		return fmt.Errorf("converting message to bytes: %v", err)
	}

	_, err = conn.WriteTo(data, addr)
	if err != nil {
		return fmt.Errorf("writing message: %v", err)
	}

	return nil
}

func ReadFrom(conn net.PacketConn, m message.Message) (net.Addr, error) {
	data := make([]byte, 1024)
	n, addr, err := conn.ReadFrom(data)
	if err != nil {
		return nil, fmt.Errorf("reading message: %v", err)
	}

	err = message.Parse(m, data[:n])
	if err != nil {
		return nil, fmt.Errorf("parsing message: %v", err)
	}

	return addr, nil
}