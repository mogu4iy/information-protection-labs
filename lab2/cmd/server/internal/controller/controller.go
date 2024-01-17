package controller

import (
	"errors"
	"fmt"
	"lab2/cmd/server/auth"
	"lab2/cmd/server/store/block"
	"lab2/cmd/server/store/user"
	"lab2/internal/constants"
	"net"
	"strings"
)


func checkUserExist(key []byte) error {
	data, err := user.Store.Search(key, 0, 1)
	if err != nil {
		return err
	}
	if len(data) != 1 {
		return fmt.Errorf("user %s does not exist", string(key))
	}
	return nil
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

func checkUserNotBlocked(key []byte) error {
	data, err := block.Store.Search(key, 0, 1)
	if err != nil {
		return err
	}
	if len(data) == 1 {
		return fmt.Errorf("user %s blocked", string(key))
	}
	return nil
}

func checkUserBlocked(key []byte) error {
	data, err := block.Store.Search(key, 0, 1)
	if err != nil {
		return err
	}
	if len(data) != 1 {
		return fmt.Errorf("user %s not blocked", string(key))
	}
	return nil
}

func onlyAdmin(addr string) error{
	if auth.Store[addr] != constants.ADMIN_USER {
		return errors.New("unauthorized")
	}
	return nil
}

func onlyUser(addr string, name string) error{
	if auth.Store[addr] != name {
		return errors.New("unauthorized")
	}
	return nil
}

func onlyAuthorized(addr string) (string, error) {
	name, ok := auth.Store[addr]
	if !ok {
		return "", errors.New("unauthorized")
	}
	return name, nil
}

func HandleAuth(conn net.PacketConn, addr net.Addr, data interface{}) (interface{}, error) {
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("data must be string")
	}
	authData := strings.Split(dataString,":")
	name := authData[0]
	password := authData[1]
	err := checkUserExist([]byte(name))
	if err != nil {
		return nil, err
	}
	err = checkUserNotBlocked([]byte(name))
	if err != nil {
		return nil, err
	}
	pHash, err := user.Store.Get([]byte(name))
	if err != nil {
		return nil, err
	}
	if !auth.CheckPasswordHash(password, string(pHash)) {
		return nil, fmt.Errorf("password is wrong")
	}
	auth.Store[addr.String()] = name
	if name == constants.ADMIN_USER {
		return constants.ADMIN_MODE, nil
	}
	return constants.USER_MODE, nil
}

func HandleCreateUser(conn net.PacketConn, addr net.Addr, data interface{}) (interface{}, error) {
	err := onlyAdmin(addr.String())
	if err != nil {
		return nil, err
	}
	err = checkUserExist([]byte(constants.ADMIN_USER))
	if err != nil {
		return nil, err
	}
	err = checkUserNotBlocked([]byte(constants.ADMIN_USER))
	if err != nil {
		return nil, err
	}
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("data must be string")
	}
	userData := strings.Split(dataString,":")
	name := userData[0]
	password := userData[1]
	err = checkUserNotExist([]byte(name))
	if err != nil {
		return nil, err
	}
	passwordHash, err := auth.HashPassword(password)
	if  err != nil {
		return nil, err
	}
	err = user.Store.Set([]byte(name), passwordHash, nil)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func HandleUpdatePassword(conn net.PacketConn, addr net.Addr, data interface{}) (interface{}, error) {
	name, err := onlyAuthorized(addr.String())
	if err != nil {
		return nil, err
	}
	err = checkUserExist([]byte(name))
	if err != nil {
		return nil, err
	}
	err = checkUserNotBlocked([]byte(name))
	if err != nil {
		return nil, err
	}
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("data must be string")
	}
	userData := strings.Split(dataString,":")
	oldPassword := userData[0]
	newPassword := userData[1]
	pHash, err := user.Store.Get([]byte(name))
	if err != nil {
		return nil, err
	}
	if !auth.CheckPasswordHash(oldPassword, string(pHash)) {
		return nil, fmt.Errorf("password is wrong")
	}
	newPasswordHash, err := auth.HashPassword(newPassword)
	if  err != nil {
		return nil, err
	}
	err = user.Store.Set([]byte(name), newPasswordHash, nil)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func HandleBlockUser(conn net.PacketConn, addr net.Addr, data interface{}) (interface{}, error) {
	err := onlyAdmin(addr.String())
	if err != nil {
		return nil, err
	}
	err = checkUserExist([]byte(constants.ADMIN_USER))
	if err != nil {
		return nil, err
	}
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("data must be string")
	}
	err = checkUserExist([]byte(dataString))
	if err != nil {
		return nil, err
	}
	err = checkUserNotBlocked([]byte(dataString))
	if err != nil {
		return nil, err
	}
	err = block.Store.Set([]byte(dataString), []byte{}, nil)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func HandleGetDoc(conn net.PacketConn, addr net.Addr, data interface{}) (interface{}, error) {
	return nil,  nil
}

func HandleStop(conn net.PacketConn, addr net.Addr) (interface{}, error) {
	delete(auth.Store, addr.String())
	return nil,  nil
}