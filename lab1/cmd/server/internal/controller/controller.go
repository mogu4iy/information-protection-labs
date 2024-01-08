package controller

import (
	"errors"
	"log"
	"regexp"
)

func FindUint(data interface{}) error {
	dataString, ok := data.(string)
	if !ok {
		return errors.New("data must be string")
	}
	re := regexp.MustCompile(`[1-9]\d*|0`)
	result := re.FindAllStringSubmatch(dataString, -1)
	log.Println(result)
	return nil
}

func FindStringKO(data interface{}) error{
	dataString, ok := data.(string)
	if !ok {
		return errors.New("data must be string")
	}
	re := regexp.MustCompile(`\b[a-zA-Z]*ko\b`)
	result := re.FindAllStringSubmatch(dataString, -1)
	log.Println(result)
	return nil
}