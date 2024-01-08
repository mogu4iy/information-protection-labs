package message

import (
	"encoding/json"
)
type Message struct {
	Command string
	Data interface{}
}

func (m *Message) ToBytes() ([]byte, error) {
	result, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *Message) Parse(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	return nil
}