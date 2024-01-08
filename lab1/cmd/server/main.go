package main

import (
	"encoding/binary"
	"fmt"
	"github.com/joho/godotenv"
	"lab1/cmd/server/internal/controller"
	"lab1/internal/constants"
	"lab1/internal/message"
	"log"
	"net"
	"os"
)

var (
	PORT string
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not loaded")
		err = nil
	}
	if Port, ok := os.LookupEnv("PORT"); ok {
		PORT = Port
	} else {
		log.Fatal("PORT env is absent")
	}
	
	listener, _ := net.Listen("tcp", ":"+PORT)
   for {
	   conn, err := listener.Accept()
	   if err != nil {
		   continue
	   }
	   go handleClient(conn)
   }
}

func handleClient(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("closing connection: %v", err)
		}
	}(conn)
	
	readLenBuf := make([]byte, 8)
	_, err := conn.Read(readLenBuf)
	if err != nil {
		log.Printf("reading message length: %v", err)
		return
	}
	readLen := binary.BigEndian.Uint64(readLenBuf)
	
	dataBuf := make([]byte, readLen)
	_, err = conn.Read(dataBuf)
	if err != nil {
		log.Printf("reading message: %v", err)
		return
	}

	m := message.Message{}
	err = m.Parse(dataBuf)
	if err != nil {
		log.Printf("parsing message: %v", err)
		return
	}

	err = commandRouter(m)
	if err != nil {
		log.Printf("handling message: %v", err)
		return 
	}
}

func commandRouter(message message.Message) error {
	switch message.Command {
	case  constants.CMDFindUint:
		err := controller.FindUint(message.Data)
		if err != nil {
			return err
		}
	case constants.CMDFindStringKO:
		err := controller.FindStringKO(message.Data)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("router for command %v not declined", message.Command)
	}
	return nil
}