package store

type Model interface {
	Parse(data []byte) error
	ToBytes() ([]byte, error)
}
