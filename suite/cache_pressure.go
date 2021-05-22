package suite

import (
	"context"
	"fmt"
	"go-redis/model"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	KEY_CACHE_PRESSURE   = "CachePressure"
	VALUE_CACHE_PRESSURE = "VALUE"
)

type CachePressure struct {
	Cache    model.Cache
	Parallel uint16
}

func (t *CachePressure) Execute() (result string, err error) {

	var (
		ctx    context.Context
		cancel context.CancelFunc
		wg     sync.WaitGroup
		count  int64
	)

	//	Cache being test
	err = t.Cache.SetValue(KEY_CACHE_PRESSURE, VALUE_CACHE_PRESSURE)
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

	result = strconv.FormatInt(count, 10)

	return
}

type result struct {
	count int64
	err   error
}

func execute(ctx context.Context, cache model.Cache, clientNumber int) <-chan result {
	//	Goroutine result channel
	ec := make(chan result)

	go func() {
		//	Generate tracffic on separate Goroutine
		ec <- generate(ctx, 60*1000, 500, 500, func(elapsed int64) result {
			fmt.Printf("[t:%d]: GET\n", clientNumber)
			value, err := cache.GetValue(KEY_CACHE_PRESSURE)
			if value != VALUE_CACHE_PRESSURE {
				return result{0, err}
			}
			return result{1, err}
		})
		close(ec)
	}()

	return ec
}

func generate(ctx context.Context, durationMs uint32, periodMs uint32, flexMs uint32, exe func(int64) result) (rt result) {
	var (
		duration int64         = int64(durationMs) * int64(time.Millisecond)
		period   time.Duration = time.Duration(periodMs) * time.Millisecond
		elapsed  int64         = 0
		now      int64         = time.Now().UnixNano()
		delay    time.Duration = time.Duration(rand.Int31n(int32(flexMs)))
		each     result
	)

	for elapsed < duration {
		rt.count += 1
		now = time.Now().UnixNano()
		delay = time.Duration(rand.Int31n(int32(flexMs)))

		select {
		case <-ctx.Done():
			rt.err = ctx.Err()
			break
		case <-time.After(delay):
			each = exe(elapsed)
			rt.count += each.count
			rt.err = each.err
			if rt.err != nil {
				break
			}
		}

		select {
		case <-ctx.Done():
			rt.err = ctx.Err()
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
