package message

import (
	"encoding/json"
)

type Message interface {
	ToBytes() ([]byte, error)
	Parse(data []byte) error
}

type Request struct {
	Command string
	Data interface{}
}

func (m *Request) ToBytes() ([]byte, error) {
	result, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *Request) Parse(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	return nil
}

type Response struct {
	Success bool
	Data interface{}
	Message string
}

func (m *Response) ToBytes() ([]byte, error) {
	result, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *Response) Parse(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	return nil
}