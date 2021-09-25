package repository

type KVRepository interface {
	AddValue(key, value []byte) error
	GetValue(key []byte) (string, error)
}
