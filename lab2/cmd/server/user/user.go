package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sopherapps/go-scdb/scdb"
	"lab2/cmd/server/internal/auth"
	"lab2/cmd/server/internal/store"
)

var Store *scdb.Store

func Init() (err error) {
	adminUser, err := json.Marshal(Data{Password: []byte("$2a$14$df8ZJxRqig.1G66pz8ZgtuuXzBFRsqi6BjTPOEjtS36dsUwSbLQtG"), IsBlocked: false})
	if err != nil {
		return
	}
	migrations := map[string][]byte{
		"ADMIN": adminUser,
	}
	Store, err = store.New( "db/user", migrations)
	if err != nil {
		return
	}
	return
}

func Exist(name string) error {
	data, err := Store.Search([]byte(name), 0, 1)
	if err != nil {
		return err
	}
	if len(data) != 1 {
		return fmt.Errorf("user %s does not exist", name)
	}
	return nil
}

type Data struct {
	Password []byte
	IsBlocked bool
}

type User struct {
	Name string
	Data Data
}

func (u *User) ToBytes() (data []byte, err error) {
	data, err = json.Marshal(&u.Data)
	if err != nil {
		return
	}
	return
}

func (u *User) Parse(data []byte) (err error) {
	err = json.Unmarshal(data, &u.Data)
	if err != nil {
		return
	}
	return
}

func (u *User) Create(name string, password string) error {
	passwordHash, err := auth.HashPassword(password)
	if  err != nil {
		return err
	}
	u.Name = name
	u.Data.Password = passwordHash
	u.Data.IsBlocked = false
	user, err := u.ToBytes()
	if err != nil {
		return err
	}
	err = Store.Set([]byte(name), user , nil)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Read(name string) error {
	user, err := Store.Get([]byte(name))
	if err != nil {
		return err
	}
	err = u.Parse(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UpdatePassword(name string, oldPassword string, password string) error {
	err := u.Read(name)
	if err != nil {
		return err
	}
	if !auth.CheckPasswordHash(oldPassword, u.Data.Password) {
		return errors.New("password not match")
	}
	passwordHash, err := auth.HashPassword(password)
	if  err != nil {
		return err
	}
	u.Data.Password = passwordHash
	user, err := u.ToBytes()
	if err != nil {
		return err
	}
	err = Store.Set([]byte(name), user, nil)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Block(name string) error {
	err := u.Read(name)
	if err != nil {
		return err
	}
	u.Data.IsBlocked = true
	user, err := u.ToBytes()
	if err != nil {
		return err
	}
	err = Store.Set([]byte(name), user, nil)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UnBlock(name string) error {
	err := u.Read(name)
	if err != nil {
		return err
	}
	u.Data.IsBlocked = false
	user, err := u.ToBytes()
	if err != nil {
		return err
	}
	err = Store.Set([]byte(name), user, nil)
	if err != nil {
		return err
	}
	return nil
}