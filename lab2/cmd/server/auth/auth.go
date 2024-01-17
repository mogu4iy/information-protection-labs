package auth

import "golang.org/x/crypto/bcrypt"

var Store = make(map[string]string)

func HashPassword(password string) ([]byte, error){
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
