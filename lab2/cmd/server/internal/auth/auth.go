package auth

import (
	"golang.org/x/crypto/bcrypt"
)

var store = make(map[string]string)

func Get(addr string) string{
	return store[addr]
}

func Login(addr string, name string){
	store[addr] = name
}

func Logout(addr string){
	delete(store, addr)
}

func IsLoggedIn(addr string) bool{
	_, ok := store[addr]
	return ok
}

func HashPassword(password string) ([]byte, error){
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
