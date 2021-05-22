package suite

import (
	"benchmark-redis/model"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	KEY_CACHE_PRESSURE   = "CachePressure"
	VALUE_CACHE_PRESSURE = "VALUE"
)

type CachePressure struct {
	Cache    func() model.Cache
	Parallel uint16
}

func (t *CachePressure) Execute() (count int64, err error) {

	var (
		c      model.Cache
		ctx    context.Context
		cancel context.CancelFunc
		wg     sync.WaitGroup
	)

	//	Cache being test
	c = t.Cache()
	defer c.Close()
	err = c.SetValue(KEY_CACHE_PRESSURE, VALUE_CACHE_PRESSURE)
	if err != nil {
		return
	}

	//	Init random
	rand.Seed(time.Now().UnixNano() / int64(time.Millisecond))

	//	Global cancel signal
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	//	Create parallel executions
	wg.Add(int(t.Parallel))
	for i := 0; i < int(t.Parallel); i++ {
		clientNumber := i
		go func() {
			defer wg.Done()
			//	Execute
			result := <-execute(ctx, t.Cache, clientNumber)
			count += result.count
			if result.err != nil && result.err != context.Canceled {
				cancel()
				err = result.err
			}
		}()
	}

	//	Wait for all worker finish
	wg.Wait()

	return
}

type result struct {
	count int64
	err   error
}

func execute(ctx context.Context, cache func() model.Cache, clientNumber int) <-chan result {
	//	Goroutine result channel
	ec := make(chan result)

	go func() {
		//	Use new connection
		c := cache()
		defer c.Close()

		//	Generate tracffic on separate Goroutine
		ec <- generate(ctx, 60*1000, 500, 500, func(elapsed int64) error {
			fmt.Printf("t:[%d]: GET\n", clientNumber)
			_, err := c.GetValue(KEY_CACHE_PRESSURE)
			return err
		})
		close(ec)
	}()

	return ec
}

func generate(ctx context.Context, durationMs uint32, periodMs uint32, flexMs uint32, exe func(int64) error) (result result) {
	var (
		duration int64         = int64(durationMs) * int64(time.Millisecond)
		period   time.Duration = time.Duration(periodMs) * time.Millisecond
		elapsed  int64         = 0
		now      int64         = time.Now().UnixNano()
		delay    time.Duration = time.Duration(rand.Int31n(int32(flexMs)))
	)

	for elapsed < duration {
		result.count += 1
		now = time.Now().UnixNano()
		delay = time.Duration(rand.Int31n(int32(flexMs)))

		select {
		case <-ctx.Done():
			result.err = ctx.Err()
			break
		case <-time.After(delay):
			result.err = exe(elapsed)
			if result.err != nil {
				break
			}
		}

		select {
		case <-ctx.Done():
			result.err = ctx.Err()
			break
		case <-time.After(max(0, period-delay)):
			elapsed += time.Now().UnixNano() - now
		}
	}
	return
}

func max(a time.Duration, b time.Duration) time.Duration {
	if a < b {
		return b
	}
	return a
}
