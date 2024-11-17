package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RespWithConn struct {
	Resp string
	Conn redis.UniversalClient
}

var conn redis.UniversalClient

//goland:noinspection ALL
func PoolHandler(ctx context.Context, req Req) ([]byte, error) {
	addrs := strings.Split(os.Getenv("REDIS_ADDRS"), ",")
	password := os.Getenv("REDIS_PASSWORD")

	randKey := fmt.Sprintf("key%d", rand.Intn(1000))

	if conn != nil {
		result, err := conn.Get(ctx, randKey).Result()
		conn.Close()
		if err == nil {
			return []byte(result), nil
		}
	}

	for i, addr := range addrs {
		addrs[i] = strings.TrimSpace(addr) + ":6379"
	}
	responseChan := make(chan RespWithConn, 3)
	cCtx, cancel := context.WithCancel(ctx)
	doneChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(3)
	for _, addr := range addrs {
		go func(addr string) {
			defer wg.Done()
			conn := redis.NewUniversalClient(
				&redis.UniversalOptions{
					Addrs:    []string{addr},
					Username: "default",
					Password: password,
					ReadOnly: true,
				},
			)
			defer conn.Close()
			result, err := conn.Get(cCtx, randKey).Result()
			if err == nil {
				responseChan <- RespWithConn{Resp: result, Conn: conn}
				doneChan <- struct{}{}
			}
		}(addr)

	}
	go func() {
		wg.Wait()
		close(responseChan)
		close(doneChan)
	}()

	<-doneChan
	cancel()

	select {
	case result := <-responseChan:
		conn = result.Conn
		return []byte(result.Resp), nil
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout waiting for response")
	}
}
