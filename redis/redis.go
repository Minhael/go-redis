package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisConfig struct {
	Ctx        context.Context
	ConnString string
}

func (c *RedisConfig) Connect() *Redis {
	return &Redis{
		ctx: c.Ctx,
		rdb: redis.NewClient(&redis.Options{
			Addr:     c.ConnString,
			Password: "",
			DB:       0,
		}),
	}
}

type Redis struct {
	ctx context.Context
	rdb *redis.Client
}

func (r *Redis) SetValue(key string, value string) error {
	return r.rdb.Set(r.ctx, key, value, 0).Err()
}

func (r *Redis) GetValue(key string) (value string, err error) {
	value, err = r.rdb.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", &NotExistError{key: key}
	}
	return
}

func (r *Redis) Close() error {
	return r.rdb.Close()
}

func (r *Redis) KeySet(pattern string, pageSize int) ([]string, error) {
	return nil, errors.New("Unimplemented")
}

func (r *Redis) Remove(keys ...string) (uint64, error) {
	return 0, errors.New("Unimplemented")
}

func (r *Redis) Size() (uint64, error) {
	return 0, errors.New("Unimplemented")
}

func (r *Redis) Clear() error {
	return errors.New("Unimplemented")
}
