package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Response struct {
	I   int
	Err error
}

func query(url string, i int, c chan Response) {
	rpcClient, err := rpc.DialHTTP(url)

	if err != nil {
		log.Printf("[%d] Failed to create a RPC client: %v", i, err)
		c <- Response{i, err}

	}

	client := ethclient.NewClient(rpcClient)

	hash := "0x75d714f13cad3b57aa240ae1f3a2a91873c994b622e582a7e19a8757d157f299"
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()
	_, err = client.TransactionReceipt(ctx, common.HexToHash(hash))

	if err != nil {
		log.Printf("[%d] Failed to query : %v", i, err)
		c <- Response{i, nil}
	}

	// log.Printf("[%d] Response %v", i, transactionReceipt)
	c <- Response{i, nil}
}

func main() {
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	url := fmt.Sprintf("https://%s:%s@%s", name, password, host)
	log.Printf("Using %s\n", url)

	c := make(chan Response)
	count := 100000
	concurrency := 100
	active := 0

	for {
		if active < concurrency && count > 0 {
			go query(url, count, c)
			count -= 1
			active += 1
		} else if active > 0 {
			response := <-c
			if response.Err != nil {
				log.Printf("[%d] %v", response.I, response.Err)
			}
			active -= 1
		} else {
			log.Printf("Request finished")
			break
		}
	}
}
