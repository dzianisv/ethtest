package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Response struct {
	I         int
	Err       error
	DelayMili int64
}

func query(url string, i int, c chan Response) (*types.Receipt, error) {
	hash := "0x75d714f13cad3b57aa240ae1f3a2a91873c994b622e582a7e19a8757d157f299"
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	start_t := time.Now()
	receipt, err := goTransactionReceipt(url, common.HexToHash(hash), ctx)
	delay_ms := time.Now().UnixMilli() - start_t.UnixMilli()

	if err != nil {
		c <- Response{i, err, delay_ms}
		return nil, err
	}

	c <- Response{i, nil, delay_ms}
	return receipt, nil
}

func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func main() {
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	url := fmt.Sprintf("https://%s:%s@%s", name, password, host)
	log.Printf("Using %s\n", url)

	c := make(chan Response)
	request_n := 100000
	count := request_n
	concurrency := 100
	active_n := 0
	errors_n := 0
	latency := make([]int64, request_n)

	for {
		if active_n < concurrency && count > 0 {
			count -= 1
			active_n += 1
			go query(url, count, c)
		} else if active_n > 0 {
			response := <-c
			latency[response.I] = response.DelayMili

			if response.Err != nil {
				errors_n += 1
				log.Printf("[%d] Error: %s; active: %d\n", response.I, response.Err, active_n)
			}
			active_n -= 1
		} else {
			break
		}
	}

	max_delay_ms := latency[0]
	sum_ms := int64(0)

	for i, delay_ms := range latency {
		max_delay_ms = Max(delay_ms, max_delay_ms)
		sum_ms += delay_ms
		log.Printf("%d\t%d\n", i, delay_ms)
	}

	avg_ms := sum_ms / int64(request_n)
	log.Printf("Max delay %d ms", max_delay_ms)
	log.Printf("Average delay %d ms", avg_ms)
	log.Printf("%d/%d failed\n", errors_n, request_n)
}
