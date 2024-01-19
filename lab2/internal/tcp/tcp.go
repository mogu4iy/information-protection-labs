package tcp

import (
	"fmt"
	"lab2/internal/kdc"
	"lab2/internal/message"
	"net"
)

func Write(conn net.Conn, m message.Message, sessionKey string) error{
	data, err := m.ToBytes()
	if err != nil {
		return fmt.Errorf("converting message to bytes: %v", err)
	}
	
	if sessionKey != "" {
		data, err = kdc.Encrypt([]byte(sessionKey), data)
		if err != nil {
			return err
		}
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("writing message: %v", err)
	}

	return nil
}

func Read(conn net.Conn, m message.Message, sessionKey string) error {
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		return fmt.Errorf("reading message: %v", err)
	}
	
	data = data[:n]
	
	if sessionKey != "" {
		data, err = kdc.Decrypt([]byte(sessionKey), data)
		if err != nil {
			return err
		}
	}
	
	err = m.Parse(data)
	if err != nil {
		return fmt.Errorf("parsing message: %v", err)
	}

	return nil
}

func Request(conn net.Conn, request message.Message, response message.Message, sessionKey string) (err error) {
	err = Write(conn, request, sessionKey)
	if err != nil {
		return
	}
	err = Read(conn, response, sessionKey)
	if err != nil {
		return
	}
	return 
}
