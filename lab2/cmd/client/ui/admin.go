package ui

import (
	"fmt"
	"lab2/cmd/client/controller"
	"lab2/internal/term"
	"log"
	"net"
	"time"
)

func Admin(conn net.Conn){
	for {
		term.Clear()
		fmt.Print("Menu:\n1. Change password\n2. Create user\n3. Block user\n4. Logout\nEnter your choice: ")
		choice := term.ReadInput()
		log.Println(choice)
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
			err := controller.CreateUser(conn)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(3 * time.Second)
			continue
		case "3":
			err := controller.BlockUser(conn)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(3 * time.Second)
			continue
		case "4":
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}