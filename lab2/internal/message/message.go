package message

import (
	"encoding/json"
)

type Message interface {
	Parse(data []byte) error
	ToBytes() ([]byte, error)
	mustImplementMessage()
}

type Request struct {
	Command string
	Data interface{}
}

func (m *Request) ToBytes() (data []byte, err error) {
	data, err = json.Marshal(m)
	if err != nil {
		return
	}
	return data, nil
}

func (m *Request) Parse(data []byte) (err error) {
	err = json.Unmarshal(data, m)
	if err != nil {
		return
	}
	return
}

func (*Request) mustImplementMessage(){}

type Response struct {
	Success bool
	Data interface{}
}

func (m *Response) ToBytes() (data []byte, err error) {
	data, err = json.Marshal(m)
	if err != nil {
		return
	}
	return data, nil
}

func (m *Response) Parse(data []byte) (err error) {
	err = json.Unmarshal(data, m)
	if err != nil {
		return
	}
	return
}

func (*Response) mustImplementMessage(){}