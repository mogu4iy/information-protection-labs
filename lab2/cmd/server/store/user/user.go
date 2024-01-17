package user

import (
	"github.com/sopherapps/go-scdb/scdb"
	"lab2/cmd/server/internal/store"
)

var Store *scdb.Store

func Init() (err error) {
	migrations := map[string][]byte{
		"ADMIN":	[]byte("$2a$14$df8ZJxRqig.1G66pz8ZgtuuXzBFRsqi6BjTPOEjtS36dsUwSbLQtG"),
	}
	Store, err = store.Init( "db/user", migrations)
	if err != nil {
		return
	}
	return
}