package block

import (
	"github.com/sopherapps/go-scdb/scdb"
	"lab2/cmd/server/internal/store"
)

var Store *scdb.Store

func Init() (err error) {
	migrations := map[string][]byte{}
	Store, err = store.Init( "db/block", migrations)
	if err != nil {
		return
	}
	return
}