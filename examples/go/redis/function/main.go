package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Req struct {
	Cmd string `json:"cmd"`
}

//goland:noinspection ALL
func PlainHandler(ctx context.Context, req Req) ([]byte, error) {
	addrs := strings.Split(os.Getenv("REDIS_ADDRS"), ",")
	password := os.Getenv("REDIS_PASSWORD")
	master := os.Getenv("REDIS_MASTER")

	for i, addr := range addrs {
		addrs[i] = strings.TrimSpace(addr) + ":26379"
	}

	conn := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs:      addrs,
			MasterName: master,
			Password:   password,
			ReadOnly:   req.Cmd != "seed",
		},
	)
	if req.Cmd == "seed" {
		seedData(ctx, conn)
		return []byte("seeded"), nil
	}

	randKey := fmt.Sprintf("key%d", rand.Intn(1000))

	result, err := conn.Get(ctx, randKey).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	_ = conn.Close()
	return []byte(result), nil
}

func seedData(ctx context.Context, conn redis.UniversalClient) {
	for i := 0; i < 1000; i++ {
		err := conn.Set(ctx, fmt.Sprintf("key%d", i), rand.Int(), 0).Err()
		if err != nil {
			panic(err)
		}
	}
}
