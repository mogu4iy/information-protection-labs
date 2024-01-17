package store

import (
	"github.com/sopherapps/go-scdb/scdb"
	"log"
)


func Init(path string, migrations map[string][]byte) (s *scdb.Store, err error) {
	var maxKeys uint64 = 1_000_000
	var redundantBlocks uint16 = 1
	var poolCapacity uint64 = 10
	var compactionInterval uint32 = 1_800
	
	s, err = scdb.New(
		path,
		&maxKeys,
		&redundantBlocks,
		&poolCapacity,
		&compactionInterval,
		true)

	for n, p := range migrations{
		kvs, err := s.Search([]byte(n), 0, 1)
		if err != nil {
			log.Fatalf("searching: %s", err)
		}
		if len(kvs) == 0 {
			err = s.Set([]byte(n), p, nil)
			if err != nil {
				log.Fatalf("setting: %s", err)
			}
		}
	}

	return
}