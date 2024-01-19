package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"lab2/cmd/server/controller"
	"lab2/cmd/server/internal/auth"
	kdcs "lab2/cmd/server/kdc"
	"lab2/cmd/server/user"
	"lab2/internal/constants"
	"lab2/internal/message"
	"lab2/internal/tcp"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	port string
)

func main() {
	err := godotenv.Load(".env.server")
	if err != nil {
		log.Println(".env file not loaded")
		err = nil
	}
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	} else {
		log.Fatal("PORT env is absent")
	}
	if s, ok := os.LookupEnv("ID_SERVER"); ok {
		kdcs.Service.ID, err = strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("ID_SERVER env is absent")
	}
	if s, ok := os.LookupEnv("MASTER_KEY"); ok {
		kdcs.Service.MasterKey = s
	} else {
		log.Fatal("MASTER_KEY env is absent")
	}

	err = user.Init()
	if err != nil {
		log.Fatalf("error initing store: %s", err)
	}
	defer func() {
		_ = user.Store.Close()
	}()
	
	TCPAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatalf("resolving TCP address: %s", err)
	}
	
	listener, err := net.ListenTCP("tcp", TCPAddr)
	if err != nil {
		log.Fatalf("creating listener: %s", err)
	}
	defer func(listener *net.TCPListener) {
		listener.Close()
	}(listener)

	log.Println("TCP server is listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		auth.Logout(conn.RemoteAddr().String())
		conn.Close()
	}(conn)
	for {
		m := &message.Request{}
		sessionKey := kdcs.Service.Users[conn.RemoteAddr().String()].SessionKey
		err := tcp.Read(conn, m, sessionKey)
		if err != nil {
			log.Println("reading: ",err)
			return
		}
		d, err := router(conn, m)
		data := &message.Response{
			Data: d,
		}
		if err != nil {
			data.Success = false
			data.Data = err.Error()
		} else {
			data.Success = true
		}
		err = tcp.Write(conn, data, sessionKey)
		if err != nil {
			return
		}
	}
}

func router(conn net.Conn, m *message.Request) (interface{}, error) {
	switch m.Command {
	case constants.HNDSHKM:
		return controller.HandleHandshake(conn, m.Data)
	case constants.AuthM:
		return controller.HandleAuth(conn, m.Data)
	case constants.CreateUserM:
		return controller.HandleCreateUser(conn, m.Data)
	case constants.UpdatePasswordM:
		return controller.HandleUpdatePassword(conn, m.Data)
	case constants.BlockUserM:
		return controller.HandleBlockUser(conn, m.Data)
//	case constants.GetDocM:
//		return controller.HandleGetDoc(conn, m.Data)
	default:
		return nil, fmt.Errorf("method not allowed")
	}
}