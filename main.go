package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Response struct {
	I         int
	Err       error
	DelayMili int64
}

func query(url string, i int, c chan Response, disableHttp2 bool, apiMethod string) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	start_t := time.Now()
	var err error = nil

	if apiMethod == "web3_clientVersion" {
		_, err = goClientVersion(url, ctx, disableHttp2)
	} else if apiMethod == "eth_transactionReceipt" {
		_, err = goTransactionReceipt(url, common.HexToHash("0x75d714f13cad3b57aa240ae1f3a2a91873c994b622e582a7e19a8757d157f299"), ctx, disableHttp2)
	}

	delay_ms := time.Now().UnixMilli() - start_t.UnixMilli()

	if err != nil {
		c <- Response{i, err, delay_ms}
		return
	}

	c <- Response{i, nil, delay_ms}
	return
}

func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

var validMethods = []string{"web3_clientVersion", "eth_transactionReceipt"}

func isSupportedMethod(method string) bool {
	for _, m := range validMethods {
		if m == method {
			return true
		}
	}

	return false
}

func main() {
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	url := fmt.Sprintf("https://%s:%s@%s", name, password, host)

	requestFlag := flag.Int("n", 1000, "number of requests")
	concurrencyFlag := flag.Int("c", 100, "number of concurrent requetss")
	disableHttp2Flag := flag.Bool("http1", false, "disable http/2")
	apiMethodFlag := flag.String("m", "web3_clientVersion", "JSON-RPC method")

	flag.Parse()

	if !isSupportedMethod(*apiMethodFlag) {
		log.Fatalf("Invalid JSON-RPC method: %s", *apiMethodFlag)
	}

	concurrency := *concurrencyFlag
	request_n := *requestFlag

	log.Printf("Using %s\n", url)
	if *disableHttp2Flag {
		log.Printf("Disable http/2")
	}

	c := make(chan Response)
	count := request_n
	active_n := 0
	errors_n := 0
	latency := make([]int64, request_n)

	for {
		if active_n < concurrency && count > 0 {
			count -= 1
			active_n += 1
			go query(url, count, c, *disableHttp2Flag, *apiMethodFlag)
		} else if active_n > 0 {
			response := <-c
			latency[response.I] = response.DelayMili

			if response.Err != nil {
				errors_n += 1
				log.Printf("[%d] Error: %s", response.I, response.Err)
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
