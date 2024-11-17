package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

type RespWithDur struct {
	Resp string
	Dur  Duration
	Addr string
}

type Response struct {
	NodeResponses []RespWithDur
	Total         Duration
}

//goland:noinspection ALL
func Handler(ctx context.Context) ([]byte, error) {
	funcStart := time.Now()
	addrs := strings.Split(os.Getenv("REDIS_ADDRS"), ",")
	password := os.Getenv("REDIS_PASSWORD")
	randKey := fmt.Sprintf("key%d", rand.Intn(1000))

	for i, addr := range addrs {
		addrs[i] = strings.TrimSpace(addr) + ":6379"
	}

	responseChan := make(chan RespWithDur, 3)
	doneChan := make(chan struct{})

	cCtx, cancel := context.WithCancel(ctx)

	for _, addr := range addrs {
		go func(addr string) {
			start := time.Now()
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
				responseChan <- RespWithDur{Resp: result, Dur: Duration{time.Since(start)}, Addr: addr}
				doneChan <- struct{}{}
			}
		}(addr)
	}

	<-doneChan
	cancel()
	close(responseChan)
	close(doneChan)

	var res []RespWithDur
	for resp := range responseChan {
		res = append(res, resp)
	}

	response := Response{
		NodeResponses: res,
		Total:         Duration{time.Since(funcStart)},
	}

	return json.Marshal(response)
}
