package model

type Set interface {
	SetValue(key string, value string) error
	GetValue(key string) (string, error)
	KeySet(pattern string, pageSize int) ([]string, error)
	Remove(keys ...string) (uint64, error)
	Size() (uint64, error)
	Clear() error
}
