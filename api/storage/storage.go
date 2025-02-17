package storage

// Storage defines a contract for every implementation detail of storage handling
type Storage interface {
	Store(string, []byte) (string, error)
	Get(string) ([]byte, error)
	Delete(string) error
}
