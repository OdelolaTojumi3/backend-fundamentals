package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Author struct {
	Name string `json: "name"`
	Age  int    `json: "age"`
}

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	json, err := json.Marshal(Author{Name: "Elliot", Age: 25})
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()

	err = client.Set(ctx, "id1234", json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	val, err := client.Get(context.Background(), "id1234").Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(val)
}
