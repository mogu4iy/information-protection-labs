package controller

import (
	"encoding/hex"
	"errors"
	"fmt"
	"lab2/cmd/server/internal/auth"
	kdcs "lab2/cmd/server/kdc"
	"lab2/cmd/server/user"
	"lab2/internal/constants"
	"lab2/internal/kdc"
	"log"
	"net"
	"strconv"
	"strings"
)


func onlyExist(name string) error {
	return user.Exist(name)
}

func checkUserNotExist(key []byte) error {
	data, err := user.Store.Search(key, 0, 1)
	if err != nil {
		return err
	}
	if len(data) == 1 {
		return fmt.Errorf("user %s exist", string(key))
	}
	return nil
}

func onlyNotBlocked(name string) error {
	data := user.User{}
	err := data.Read(name)
	if err != nil {
		return err
	}
	if data.Data.IsBlocked {
		return errors.New("user blocked")
	}
	return nil
}

func onlyAdmin(addr string) error{
	if auth.Get(addr) != constants.AdminUser {
		return errors.New("unauthorized")
	}
	return nil
}

func onlyAuth(addr string) error {
	ok := auth.IsLoggedIn(addr)
	if !ok {
		return errors.New("unauthorized")
	}
	return nil
}

func HandleHandshake(conn net.Conn, data interface{}) (interface{}, error) {
	addr := conn.RemoteAddr().String()
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("data type is wrong")
	}
	dataBytes, err := hex.DecodeString(dataString)
	if err != nil {
		return nil, err
	}
	
	requestData := strings.Split(string(dataBytes),":")
	serverData := requestData[0]
	sessionData := requestData[1]
	
	serverDataDecrypted, err := kdc.Decrypt([]byte(kdcs.Service.MasterKey), []byte(fmt.Sprintf("%s", serverData)))
	if err != nil {
		return nil, err
	}
	serverDataDecryptedStirng := string(serverDataDecrypted)
	serverDataArray := strings.Split(serverDataDecryptedStirng, ":")
	userID, err := strconv.Atoi(serverDataArray[1])
	if err != nil {
		return nil, err
	}
	
	kdcs.Service.Users[addr] = kdc.User{ID: userID, SessionKey: serverDataArray[0]}
	
	sessionDataDecrypted, err := kdc.Decrypt([]byte(serverDataArray[0]), []byte(sessionData))
	if err != nil {
		return nil, err
	}
	sessionDataDecryptedStirng := string(sessionDataDecrypted)
	r2, err := strconv.Atoi(sessionDataDecryptedStirng)
	if err != nil {
		return nil, err
	}
	r2Func := kdc.RandFunc(r2)
	clientSessionData, err := kdc.Encrypt([]byte(serverDataArray[0]), []byte(fmt.Sprintf("%d", r2Func)))
	if err != nil {
		return nil, err
	}

	return hex.EncodeToString(clientSessionData), nil
}

func HandleAuth(conn net.Conn, data interface{}) (interface{}, error) {
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("data type is wrong")
	}
	
	requestData := strings.Split(dataString,":")
	name := requestData[0]
	password := requestData[1]
	
	err := onlyNotBlocked(name)
	if err != nil {
		return nil, err
	}
	err = onlyExist(name)
	if err != nil {
		return nil, err
	}
	
	userModel := user.User{}
	err = userModel.Read(name)
	if err != nil {
		return nil, err
	}
	if !auth.CheckPasswordHash(password, userModel.Data.Password) {
		return nil, fmt.Errorf("password is wrong")
	}
	
	auth.Login(conn.RemoteAddr().String(), name)
	
	return "authenticated", nil
}

func HandleCreateUser(conn net.Conn, data interface{}) (interface{}, error) {
	err := onlyAdmin(conn.RemoteAddr().String())
	if err != nil {
		return nil, err
	}
	err = onlyExist(constants.AdminUser)
	if err != nil {
		return nil, err
	}
	err = onlyNotBlocked(constants.AdminUser)
	if err != nil {
		return nil, err
	}
	
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("data type is wrong")
	}
	
	requestData := strings.Split(dataString,":")
	name := requestData[0]
	password := requestData[1]
	
	err = onlyExist(name)
	if err == nil {
		return nil, errors.New("user already exist")
	}
	
	userData := &user.User{}
	err = userData.Create(name, password)
	if err != nil {
		return nil, err
	}
	
	return "user created", nil
}

func HandleUpdatePassword(conn net.Conn, data interface{}) (interface{}, error) {
	err := onlyAuth(conn.RemoteAddr().String())
	if err != nil {
		return nil, err
	}
	name := auth.Get(conn.RemoteAddr().String())
	err = onlyExist(name)
	if err != nil {
		return nil, err
	}
	err = onlyNotBlocked(name)
	if err != nil {
		return nil, err
	}
	
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("data must be string")
	}
	
	requestData := strings.Split(dataString,":")
	oldPassword := requestData[0]
	password := requestData[1]
	
	userData := &user.User{}
	err = userData.UpdatePassword(name, oldPassword, password)
	if err != nil {
		return nil, err
	}
	
	return "password updated", nil
}

func HandleBlockUser(conn net.Conn, data interface{}) (interface{}, error) {
	err := onlyAdmin(conn.RemoteAddr().String())
	if err != nil {
		return nil, err
	}
	err = onlyExist(constants.AdminUser)
	if err != nil {
		return nil, err
	}
	err = onlyNotBlocked(constants.AdminUser)
	if err != nil {
		return nil, err
	}
	
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("data must be string")
	}
	
	err = onlyExist(dataString)
	if err != nil {
		return nil, err
	}
	err = onlyNotBlocked(dataString)
	if err == nil {
		return nil, errors.New("user already blocked")
	}
	
	userData := &user.User{}
	err = userData.Block(dataString)
	if err != nil {
		return nil, err
	}
	
	return "user blocked", nil
}

func HandleGetDoc(conn net.Conn, data interface{}) (interface{}, error) {
	return nil,  nil
}