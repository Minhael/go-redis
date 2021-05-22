package model

type Cache interface {
	SetValue(key string, value string) error
	GetValue(key string) (string, error)
	Close() error
}
