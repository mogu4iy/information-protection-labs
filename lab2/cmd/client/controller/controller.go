package controller

import (
	"encoding/hex"
	"errors"
	"fmt"
	kdcs "lab2/cmd/client/kdc"
	"lab2/internal/constants"
	"lab2/internal/kdc"
	"lab2/internal/message"
	"lab2/internal/tcp"
	"lab2/internal/term"
	"net"
	"strconv"
)

func Handshake(conn net.Conn, serverData string) error {
	r2 := kdc.GenerateRandomNumber(0,100)
	sessionData, err := kdc.Encrypt([]byte(kdcs.Service.SessoinKey), []byte(fmt.Sprintf("%d", r2)))
	if err != nil {
		return err
	}
	
	request := &message.Request{
		Command: constants.HNDSHKM,
		Data: hex.EncodeToString([]byte(fmt.Sprintf("%s:::%s", serverData, sessionData))),
	}
	response := &message.Response{}

	err = tcp.Request(conn, request, response, "")
	if err != nil {
		return err
	}
	
	if !response.Success {
		return fmt.Errorf("%s", response.Data)
	}
	
	responseByte, err := hex.DecodeString(response.Data.(string))
	if err != nil {
		return err
	}
	
	r2Func, err := kdc.Decrypt([]byte(kdcs.Service.SessoinKey), responseByte)
	if err != nil {
		return err
	}
	r2FuncInt, err := strconv.Atoi(string(r2Func))
	if err != nil {
		return err
	}
	if r2FuncInt != kdc.RandFunc(r2) {
		return errors.New("f(r2) not match")
	}
	
	return nil
}

func Auth(conn net.Conn, name string, password string) error {
	request := &message.Request{
		Command: constants.AuthM,
		Data: fmt.Sprintf("%v:%v", name, password),
	}
	response := &message.Response{}
	
	err := tcp.Request(conn, request, response, kdcs.Service.SessoinKey)
	if err != nil {
		return err
	}
	
	if !response.Success {
		return fmt.Errorf("%s", response.Data)
	} else {
		fmt.Println(response.Data)
	}

	return nil
}

func UpdatePassword(conn net.Conn) error {
	oldPassword := term.ReadVar("old password", false)
	password := term.ReadVar("password", true)
	fmt.Println()
	
	request := &message.Request{
		Command: constants.UpdatePasswordM,
		Data: fmt.Sprintf("%v:%v", oldPassword, password),
	}
	response := &message.Response{}
	
	err := tcp.Request(conn, request, response, kdcs.Service.SessoinKey)
	if err != nil {
		return err
	}
	
	if !response.Success {
		return fmt.Errorf("%s", response.Data)
	} else {
		fmt.Println(response.Data)
	}
	
	return nil
}

func CreateUser(conn net.Conn) error {
	name := term.ReadVar("name", false)
	password := term.ReadVar("password", true)
	fmt.Println()
	
	request := &message.Request{
		Command: constants.CreateUserM,
		Data: fmt.Sprintf("%v:%v", name, password),
	}
	response := &message.Response{}
	
	err := tcp.Request(conn, request, response, kdcs.Service.SessoinKey)
	if err != nil {
		return err
	}
	
	if !response.Success {
		return fmt.Errorf("%s", response.Data)
	} else {
		fmt.Println(response.Data)
	}
	
	return nil
}

func BlockUser(conn net.Conn) error {
	name := term.ReadVar("name", false)
	fmt.Println()
	
	request := &message.Request{
		Command: constants.BlockUserM,
		Data: name,
	}
	response := &message.Response{}
	
	err := tcp.Request(conn, request, response, kdcs.Service.SessoinKey)
	if err != nil {
		return err
	}
	
	if !response.Success {
		return fmt.Errorf("%s", response.Data)
	} else {
		fmt.Println(response.Data)
	}
	
	return nil
}