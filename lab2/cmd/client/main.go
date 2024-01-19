package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"lab2/cmd/client/controller"
	kdcs "lab2/cmd/client/kdc"
	"lab2/cmd/client/ui"
	"lab2/internal/constants"
	"lab2/internal/kdc"
	"lab2/internal/term"
	"lab2/internal/udp"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	nameFlag string
	passwordFlag string
	serverAddr string
	kdcAddr string
	IDServer int
)

func main() {
	err := godotenv.Load(".env.client")
	if err != nil {
		log.Println(".env file not loaded")
		err = nil
	}
	if s, ok := os.LookupEnv("SERVER_ADDR"); ok {
		serverAddr = s
	} else {
		log.Fatal("SERVER_ADDR env is absent")
	}
	if s, ok := os.LookupEnv("KDC_ADDR"); ok {
		kdcAddr = s
	} else {
		log.Fatal("KDC_ADDR env is absent")
	}
	if s, ok := os.LookupEnv("ID_SERVER"); ok {
		IDServer, err = strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("ID_SERVER env is absent")
	}
	if s, ok := os.LookupEnv("ID_CLIENT"); ok {
		kdcs.Service.ID, err = strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("ID_CLIENT env is absent")
	}
	if s, ok := os.LookupEnv("MASTER_KEY"); ok {
		kdcs.Service.MasterKey = s
	} else {
		log.Fatal("MASTER_KEY env is absent")
	}

	sessionKeyData, err := getSessionKey()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal("connecting to server: ", err)
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatalf("closing connection: %s", err)
		}
	}(conn)

	err = controller.Handshake(conn, sessionKeyData)
	if err != nil {
		log.Fatal(err)
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

	err = controller.Auth(conn, nameFlag, passwordFlag)
	if err != nil{
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	if nameFlag == constants.AdminUser {
		ui.Admin(conn)

	} else {
		ui.User(conn)
	}
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

func getSessionKey() (string, error) {
	conn, err := net.Dial("udp", kdcAddr)
	if err != nil {
		log.Fatal("connecting to KDC server: ", err)
	}

	r1 := kdc.GenerateRandomNumber(0, 100)
	kdcData, err := kdc.Encrypt([]byte(kdcs.Service.MasterKey), []byte(fmt.Sprintf("%d:%d", r1, IDServer)))
	if err != nil {
		return "", err
	}

	request := []byte(fmt.Sprintf("%d:%s", kdcs.Service.ID, kdcData))
	err = udp.Write(conn, request)
	if err != nil {
		return "", err
	}
	response, err := udp.Read(conn)
	if err != nil {
		return "", err
	}

	sessionKeyData, err := kdc.Decrypt([]byte(kdcs.Service.MasterKey), response)
	if err != nil {
		return "", err
	}
	sessionKeyDataString := string(sessionKeyData)
	sessionKeyDataArray := strings.Split(sessionKeyDataString, ":")

	r1Func, err := strconv.Atoi(sessionKeyDataArray[0])
	if err != nil {
		return "", err
	}
	if kdc.RandFunc(r1) != r1Func {
		return "", errors.New("f(r1) not match")
	}
	kdcs.Service.SessoinKey = sessionKeyDataArray[1]
	IDB, err := strconv.Atoi(sessionKeyDataArray[2])
	if err != nil {
		return "", err
	}
	if IDServer != IDB {
		return "", errors.New("IDServer not match")
	}

	return sessionKeyDataArray[3], nil
}