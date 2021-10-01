package repository

type KVRepository interface {
	AddValue(key, value []byte) error
	GetValue(key []byte) (string, error)
	Close() error
}
type SQLRepository interface {
	GetAllValues() []*SqlitePerson
	Close() error
}
