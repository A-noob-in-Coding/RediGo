package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"os"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)
var Origin string
var Port string // default port
func handleReq(v http.ResponseWriter, r *http.Request) {
	requestURL := Origin
	header := v.Header()
	reqHeaders := r.Header.Get("Cache-Required")
	if r.URL.Path == "/redigo/clear" {
		err := rdb.FlushDB(ctx).Err()
		if err != nil {
			panic(err)
		}
		io.WriteString(v, "Cache Cleared Successfully\n")
		return
	}

	if r.URL.Path == "/" {
		requestURL += r.URL.Path[1:]
	} else {
		requestURL += r.URL.Path
	}

	result, err := rdb.Get(ctx, requestURL).Result()
	if err == nil && reqHeaders != "no" {

		fmt.Printf("Cache Hit\n")
		header.Set("X-Cache", "Hit")
		io.WriteString(v, result)

		return
	}
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("could not get")
	}
	if resp.StatusCode == 404 {
		io.WriteString(v, "ERROR 404, invalid api endpoint")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if reqHeaders == "no" {
		header.Set("X-Cache", "Refreshed")
		fmt.Printf("Cache Refreshed\n")
	} else {
		fmt.Printf("Cache Miss\n")
		header.Set("X-Cache", "Miss")
	}
	_, _ = io.WriteString(v, string(body))

	err = rdb.Set(ctx, requestURL, string(body), 0).Err()
	if err != nil {
		fmt.Println("Error storing in Redis:", err)
	}
}

func main() {
	var exists bool
	Port, exists = os.LookupEnv("PORT")
	if !exists {
		Port = "1234"
	}

	Origin, exists = os.LookupEnv("ORIGIN")
	if !exists || Origin == "" {
		fmt.Println("Set ORIGIN environment variable")
		return
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	http.HandleFunc("/", handleReq)

	fmt.Printf("Server running at Port %s\n", Port)

	err := http.ListenAndServe(":"+Port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
