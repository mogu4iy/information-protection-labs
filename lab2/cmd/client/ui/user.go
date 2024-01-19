package ui

import (
	"fmt"
	"lab2/cmd/client/controller"
	"lab2/internal/term"
	"log"
	"net"
	"time"
)

func User(conn net.Conn) {
	for {
		term.Clear()
		fmt.Print("Menu:\n1. Change password\n2. Get document\n3. Logout\nEnter your choice: ")
		choice := term.ReadInput()
		term.Clear()
		switch choice {
		case "1":
			err := controller.UpdatePassword(conn)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(3 * time.Second)
			continue
		case "2":
			fmt.Println("Unsupported. Try Again.")
			continue
		case "3":
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}