package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"lab2/cmd/server/internal/controller"
	"lab2/cmd/server/store/block"
	"lab2/cmd/server/store/user"
	"lab2/internal/constants"
	"lab2/internal/message"
	"lab2/internal/udp"
	"log"
	"net"
	"os"
)

var (
	Port string
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not loaded")
		err = nil
	}
	if p, ok := os.LookupEnv("PORT"); ok {
		Port = p
	} else {
		log.Fatal("PORT env is absent")
	}

	err = user.Init()
	if err != nil {
		log.Fatalf("error initing store: %s", err)
	}
	defer func() {
		_ = user.Store.Close()
	}()
	
	err = block.Init()
	if err != nil {
		log.Fatalf("error initing store: %s", err)
	}
	defer func() {
		_ = block.Store.Close()
	}()
	
	
	TCPAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%v", Port))
	if err != nil {
		log.Fatalf("resolving UDP address: %s", err)
	}
	
	conn, err := net.ListenUDP("udp", TCPAddr)
	if err != nil {
		log.Fatalf("listening: %s", err)
		return
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatalf("closing connection: %s", err)
		}
	}(conn)
	
	log.Println("UDP server is listening")

	for {
		m := &message.Request{}
		addr, err := udp.ReadFrom(conn, m)
		if err != nil {
			log.Println(err)
			continue
		}
		go handleRequest(conn, addr, m)
	}
}

func handleRequest(conn net.PacketConn, addr net.Addr, m *message.Request) {
	d, err := router(conn,addr, m)
	data := &message.Response{
		Data: d,
	}
	if err != nil {
		data.Success = false
		data.Message = err.Error()
	} else {
		data.Success = true
	}
	err = udp.WriteTo(conn, addr, data)
	if err != nil {
		log.Println(err)
		return 
	}
}

func router(conn net.PacketConn, addr net.Addr, m *message.Request) (interface{}, error) {
	switch m.Command {
	case constants.AUTH_M:
		return controller.HandleAuth(conn, addr, m.Data)
	case constants.CREATE_USER_M:
		return controller.HandleCreateUser(conn, addr, m.Data)
	case constants.UPDATE_PASSWORD_M:
		return controller.HandleUpdatePassword(conn, addr, m.Data)
	case constants.BLOCK_USER_M:
		return controller.HandleBlockUser(conn, addr, m.Data)
	case constants.GET_DOC_M:
		return controller.HandleGetDoc(conn, addr, m.Data)
	case constants.STOP_M:
		return controller.HandleStop(conn, addr)
	default:
		return nil, fmt.Errorf("method not allowed")
	}
}