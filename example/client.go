package main

import (
	"github.com/go-redis/redis"
	"runtime"
	"fmt"
)

func main()  {
	option := &redis.Options{
		Addr:       "127.0.0.1:6789",
		Password:   "",
		DB:         0,
		MaxRetries: 7,
		PoolSize:   20 * runtime.NumCPU(),
	}

	for i := 0; i < 10000; i++ {
		client := redis.NewClient(option)
		cmd := client.Get("name")
		data,_ := cmd.Result()
		client.Close()
		fmt.Printf("%d-%s\n", i, data)
	}

}
