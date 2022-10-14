package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Response struct {
	I   int
	Err error
}

func query(url string, i int, c chan Response) (*types.Receipt, error) {
	rpcClient, err := rpc.DialHTTP(url)

	if err != nil {
		c <- Response{i, err}
		return nil, err
	}

	client := ethclient.NewClient(rpcClient)

	hash := "0x75d714f13cad3b57aa240ae1f3a2a91873c994b622e582a7e19a8757d157f299"
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()
	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(hash))

	if err != nil {
		c <- Response{i, err}
		return nil, err
	}

	c <- Response{i, nil}
	return receipt, nil
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

	for {
		if active_n < concurrency && count > 0 {
			go query(url, count, c)
			count -= 1
			active_n += 1
		} else if active_n > 0 {
			response := <-c
			if response.Err != nil {
				errors_n += 1
				log.Printf("[%d] Error: %s; active: %d\n", response.I, response.Err, active_n)
			}
			active_n -= 1
		} else {
			log.Printf("%d/%d failed\n", errors_n, request_n)
			break
		}
	}
}
