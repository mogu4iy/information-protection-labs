package kdc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	mrand "math/rand"
)
	
var letterRunes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GenerateMasterKey() ([]byte, error) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = letterRunes[mrand.Intn(len(letterRunes))]
	}
	return key, nil
}

func GenerateRandomNumber(min, max int) int {
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	randomInt := int(binary.BigEndian.Uint32(randomBytes))
	return min + randomInt%(max-min+1)
}

//func NumberToBytes(value int) []byte {
//	bytes := make([]byte, 4)
//	binary.BigEndian.PutUint32(bytes, uint32(value))
//	if len(bytes) < 15 {
//		padding := make([]byte, 15-len(bytes))
//		bytes = append(bytes, padding...)
//	}
//	return bytes
//}

func GenerateSessionKey(key1 []byte, key2 []byte) []byte {
	sessionKey := make([]byte, 32)
	for i := 0; i < 32; i++ {
		sessionKey[i] = key1[i] ^ key2[i]
	}
	return sessionKey
}

func Encrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	return append(nonce, ciphertext...), nil
}

func Decrypt(key []byte, encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func RandFunc(r int) int {
	return r * 10
}

type User struct {
	ID int
	SessionKey string
}

type ServiceClient struct {
	ID int
	MasterKey string
	SessoinKey string
}

type ServiceService struct {
	ID int
	MasterKey string
	Users map[string]User
}