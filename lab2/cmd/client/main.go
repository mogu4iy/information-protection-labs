package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"lab2/internal/constants"
	"lab2/internal/message"
	"lab2/internal/term"
	"lab2/internal/udp"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	nameFlag string
	passwordFlag string
	ServerAddr string
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not loaded")
		err = nil
	}
	if s, ok := os.LookupEnv("SERVER_ADDR"); ok {
		ServerAddr = s
	} else {
		log.Fatal("PORT env is absent")
	}


	flag.StringVar(&nameFlag, "n", "", "user name")
	flag.StringVar(&passwordFlag, "p", "", "user password")
	flag.Parse()
	if !isFlagPassed("n") {
		nameFlag = term.ReadVar("name", false)
	}
	if !isFlagPassed("p") {
		passwordFlag = term.ReadVar("password", true)
	}
	
	conn, err := net.Dial("udp", ServerAddr)
	if err != nil {
		log.Fatalf("connecting to server: %s", err)
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatalf("closing connection: %s", err)
		}
	}(conn)
	
	authRequest := &message.Request{
		Command: constants.AUTH_M,
		Data: fmt.Sprintf("%v:%v", nameFlag, passwordFlag),
	}
	err = udp.Write(conn, authRequest)
	if err != nil {
		log.Fatal(err)
	}
	
	authResponse := &message.Response{}
	err = udp.Read(conn, authResponse)
	if err != nil {
		log.Fatal(err)
	}
	if !authResponse.Success {
		log.Fatal(authResponse.Message)
	}
	switch authResponse.Data {
	case constants.ADMIN_MODE:
		adminMode(conn)
		return
	case constants.USER_MODE:
		userMode(conn)
		return
	default:
		log.Fatalf("unsupported mode received: %v", authResponse.Data)
	}
}

func adminMode(conn net.Conn) {
	for {
		clearTerminal()
		fmt.Print("Menu:\n1. Change password\n2. Create user\n3. Block user\n4. Logout\nEnter your choice: ")
		choice := readInput()
		clearTerminal()
		switch choice {
		case "1":
			oldPassword := term.ReadVar("old password", false)
			password := term.ReadVar("password", true)
			request := &message.Request{
				Command: constants.UPDATE_PASSWORD_M,
				Data: fmt.Sprintf("%v:%v", oldPassword, password),
			}
			err := udp.Write(conn, request)
			if err != nil {
				log.Fatal(err)
			}
			response := &message.Response{}
			err = udp.Read(conn, response)
			if err != nil {
				log.Fatal(err)
			}
			if !response.Success {
				log.Println(response.Message)
			} else {
				log.Println("Success!")
			}
			time.Sleep(3 * time.Second)
			continue
		case "2":
			name := term.ReadVar("name", false)
			password := term.ReadVar("password", true)
			request := &message.Request{
				Command: constants.CREATE_USER_M,
				Data: fmt.Sprintf("%v:%v", name, password),
			}
			err := udp.Write(conn, request)
			if err != nil {
				log.Fatal(err)
			}
			response := &message.Response{}
			err = udp.Read(conn, response)
			if err != nil {
				log.Fatal(err)
			}
			if !response.Success {
				log.Println(response.Message)
			} else {
				log.Println("Success!")
			}
			time.Sleep(3 * time.Second)
			continue
		case "3":
			name := term.ReadVar("name", false)
			request := &message.Request{
				Command: constants.BLOCK_USER_M,
				Data: name,
			}
			err := udp.Write(conn, request)
			if err != nil {
				log.Fatal(err)
			}
			response := &message.Response{}
			err = udp.Read(conn, response)
			if err != nil {
				log.Fatal(err)
			}
			if !response.Success {
				log.Println(response.Message)
			} else {
				log.Println("Success!")
			}
			time.Sleep(3 * time.Second)
			continue
		case "4":
			request := &message.Request{
				Command: constants.STOP_M,
			}
			err := udp.Write(conn, request)
			if err != nil {
				log.Fatal(err)
			}
			response := &message.Response{}
			err = udp.Read(conn, response)
			if err != nil {
				log.Fatal(err)
			}
			if !response.Success {
				log.Println(response.Message)
			} else {
				log.Println("Success!")
			}
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}

func userMode(conn net.Conn) {
	for {
		clearTerminal()
		fmt.Print("Menu:\n1. Change password\n2. Get document\n3. Logout\nEnter your choice: ")
		choice := readInput()
		clearTerminal()
		switch choice {
		case "1":
			oldPassword := term.ReadVar("old password", false)
			password := term.ReadVar("password", true)
			request := &message.Request{
				Command: constants.UPDATE_PASSWORD_M,
				Data: fmt.Sprintf("%v:%v", oldPassword, password),
			}
			err := udp.Write(conn, request)
			if err != nil {
				log.Fatal(err)
			}
			response := &message.Response{}
			err = udp.Read(conn, response)
			if err != nil {
				log.Fatal(err)
			}
			if !response.Success {
				log.Println(response.Message)
			} else {
				log.Println("Success!")
			}
			time.Sleep(3 * time.Second)
			continue
		case "2":
			continue
		case "3":
			request := &message.Request{
				Command: constants.STOP_M,
			}
			err := udp.Write(conn, request)
			if err != nil {
				log.Fatal(err)
			}
			response := &message.Response{}
			err = udp.Read(conn, response)
			if err != nil {
				log.Fatal(err)
			}
			if !response.Success {
				log.Println(response.Message)
			} else {
				log.Println("Success!")
			}
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}

func clearTerminal() {
	fmt.Print("\033c")
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}