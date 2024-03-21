package db

type Database interface {
	Open(dir string) error
	Close() error
	Read(key string) ([]byte, error)
	Write(key string, value []byte) error
	Delete(key string) error
	ReadAll() (map[string][]byte, error)
}
