package repository

type KVRepository interface {
	AddValue(bucketName,key, value []byte) error
	GetValue(bucketName,key []byte) (string, error)
	Close() error
}
// type SQLRepository interface {
// 	GetAllValues() []*SqlitePerson
// 	Close() error
// }
