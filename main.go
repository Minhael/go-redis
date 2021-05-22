package main

import (
	"benchmark-redis/model"
	"benchmark-redis/redis"
	"benchmark-redis/suite"
	"context"
	"fmt"
	"time"
)

const (
	connString   = "localhost:6379"
	clientNumber = 300
)

func main() {

	//	Pressure test with multiple clients
	suite := &suite.CachePressure{
		Cache:    connect,
		Parallel: clientNumber,
	}

	now := time.Now().UnixNano()
	result, err := suite.Execute()
	elapsed := (time.Now().UnixNano() - now) / int64(time.Millisecond)
	if err != nil {
		fmt.Printf("Failed in %d ms size %d\n%s", elapsed, result, err)
	} else {
		fmt.Printf("Completed in %d ms size %d", elapsed, result)
	}
}

func connect() model.Cache {
	ctx := context.Background()
	config := &redis.RedisConfig{
		Ctx:        ctx,
		ConnString: connString,
	}
	return config.Connect()
}
