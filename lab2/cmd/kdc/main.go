package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"lab2/internal/kdc"
	"lab2/internal/udp"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID	int
	Key []byte
}

type KeyDistributionCenter struct {
	users map[int]*User
}

var KDC = &KeyDistributionCenter{
	users: map[int]*User{
		1: &User{
			ID: 1,
			Key: []byte("nr2bmdeOOYimz48GV7hua26qRMtfKEFZ"),
		},
		2: &User{
			ID: 2,
			Key: []byte("s45lbZd4hB1zhaLHofhyzh5BShRpH5Wv"),
			},
	},
}

func handleUDPConnection(conn net.PacketConn, addr net.Addr, m []byte) {
	dataString := string(m)
	dataArray := strings.Split(dataString,":::")
	IDA, err := strconv.Atoi(dataArray[0])
	if err != nil {
		return
	}
	data := dataArray[1]
	
	decryptedData, err := kdc.Decrypt(KDC.users[IDA].Key, []byte(data))
	decryptedDataString := string(decryptedData)
	decryptedDataArray := strings.Split(decryptedDataString, ":::")
	r1, err := strconv.Atoi(decryptedDataArray[0])
	if err != nil {
		return
	}
	IDB, err := strconv.Atoi(decryptedDataArray[1])
	if err != nil {
		return
	}
	
	sessionKey := kdc.GenerateSessionKey(KDC.users[IDA].Key, KDC.users[IDB].Key)
	sessionData := fmt.Sprintf("%s:::%d", sessionKey, IDA)
	log.Printf("%d, %s", IDB, KDC.users[IDB].Key)
	sessionDataEncrypted, err := kdc.Encrypt(KDC.users[IDB].Key, []byte(sessionData))
	if err != nil {
		return
	}
	log.Printf("%v \n", sessionData)
	aData := fmt.Sprintf("%d:::%s:::%d:::%s", kdc.RandFunc(r1), sessionKey, IDB, sessionDataEncrypted)
	
	aDataEncrypted, err := kdc.Encrypt(KDC.users[IDA].Key, []byte(aData))
	if err != nil {
		return
	}
	
	err = udp.WriteTo(conn, addr, aDataEncrypted)
	if err != nil {
		return 
	}
	
}

var (
	port string
)

func main() {
	err := godotenv.Load(".env.kdc")
	if err != nil {
		log.Println(".env file not loaded")
		err = nil
	}
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	} else {
		log.Fatal("PORT env is absent")
	}
	
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		fmt.Println("resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP server is running")

	for {
		m, addr, err := udp.ReadFrom(conn)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			return
		}
		go handleUDPConnection(conn, addr, m)
	}
}

//func main() {
//	serverMasterKey, err := kdc.GenerateMasterKey()
//	if err != nil {
//		log.Fatal(err)
//	}
//	clientMasterKey, err := kdc.GenerateMasterKey()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%s", serverMasterKey)
//	fmt.Println()
//	fmt.Println()
//	fmt.Println()
//	fmt.Printf("%s", clientMasterKey)
//}