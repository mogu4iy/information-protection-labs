package kdc

import "lab2/internal/kdc"

var Service = &kdc.ServiceService{
	Users: make(map[string]kdc.User),
}