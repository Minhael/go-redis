package main

import (
	"benchmark-redis/redis"
	"benchmark-redis/suite"
	"context"
	"fmt"
	"time"
)

const (
	connString   = "localhost:6379"
	clientNumber = 80
)

func main() {

	//	Go Redis library increase number of clients internally if connected clients are not enough to serve the pressure
	//	Connect to database
	ctx := context.Background()
	config := &redis.RedisConfig{
		Ctx:        ctx,
		ConnString: connString,
	}

	//	Pressure test with multiple clients
	suite := &suite.CachePressure{
		Cache:    config.Connect(),
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
