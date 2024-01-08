package main

import (
	"encoding/binary"
	"github.com/joho/godotenv"
	"lab1/internal/constants"
	"lab1/internal/message"
	"log"
   "net"
	"os"
)

var (
   ServerAddr string
)

func main() {
   err := godotenv.Load(".env")
   if err != nil {
      log.Println(".env file not loaded")
      err = nil
   }
   if serverAddr, ok := os.LookupEnv("SERVER_ADDR"); ok {
      ServerAddr = serverAddr
   } else {
      log.Fatal("SERVER_ADDR env is absent")
   }
   
   conn, _ := net.Dial("tcp", ServerAddr)
   defer func(conn net.Conn) {
	   err := conn.Close()
	   if err != nil {
         log.Printf("closing connection: %v", err)
      }  
   }(conn)
   
   m := message.Message{
      Command: constants.CMDFindStringKO,
      Data: "Zaichenko 10ko ko rabbits had been eating 100 carrots loudly, so 1.5 wolves easily found them. As a result 9 rabbits were saying 'rest in peace rabbit Zaichenko'",
   }
   b, err := m.ToBytes()
   if err != nil {
      log.Printf("converting message to bytes: %v", err)
      return
   }
   
   data := make([]byte, 8)
   binary.BigEndian.PutUint64(data, uint64(len(b)))
	data = append(data, b...)

   _, err = conn.Write(data)
	if err != nil {
      log.Printf("writing message: %v", err)
      return
	}
}