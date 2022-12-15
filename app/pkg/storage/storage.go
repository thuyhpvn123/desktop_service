package storage

type IStorage interface {
	Get([]byte) ([]byte, error)
	Put([]byte, []byte) error
	Has([]byte) bool
	Delete([]byte) error
	BatchPut([][2][]byte) error
	Close()
	GetIterator() IIterator
}
