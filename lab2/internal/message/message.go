package message

import (
	"encoding/json"
	"fmt"
	"github.com/howeyc/crc16"
)

type Message interface {
	mustImplementMessage()
}

func ToBytes(m Message) ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	// Додати обчислення CRC16 з зсувом на 7 біт
	crc := crc16.ChecksumCCITT(data) << 7
	crcBytes := []byte{byte(crc >> 8), byte(crc & 0xFF)}
	data = append(data, crcBytes...)
	return data, nil
}

func Parse(m Message, data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("недостатньо даних для перевірки CRC16")
	}
	receivedCRC := uint16(data[len(data)-2])<<8 | uint16(data[len(data)-1])
	dataWithoutCRC := data[:len(data)-2]

	// Зсув на 7 біт для обчислення CRC16
	calculatedCRC := crc16.ChecksumCCITT(dataWithoutCRC) << 7
	if receivedCRC != calculatedCRC {
		return fmt.Errorf("неправильна контрольна сума CRC16")
	}

	err := json.Unmarshal(dataWithoutCRC, m)
	if err != nil {
		return err
	}
	return nil
}

type Request struct {
	Command string
	Data interface{}
}

func (*Request) mustImplementMessage(){}

type Response struct {
	Success bool
	Data interface{}
	Message string
}

func (*Response) mustImplementMessage(){}