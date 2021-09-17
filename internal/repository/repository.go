package repository

type KVRepository interface {
	AddValue(key, value string) error
	GetValue(key string) (string, error)
}
