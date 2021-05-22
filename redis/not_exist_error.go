package redis

type NotExistError struct {
	key string
}

func (e *NotExistError) Error() string {
	return "Key not exists: " + e.key
}
