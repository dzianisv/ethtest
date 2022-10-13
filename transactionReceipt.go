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

func main() {
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	url := fmt.Sprintf("https://%s:%s@%s", name, password, host)

	log.Printf("Using %s\n", url)

	rpcClient, err := rpc.DialHTTP(url)

	if err != nil {
		log.Fatalf("Failed to create a RPC client: %v", err)
	}

	client := ethclient.NewClient(rpcClient)

	hash := "0x75d714f13cad3b57aa240ae1f3a2a91873c994b622e582a7e19a8757d157f299"
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()
	transactionReceipt, err := client.TransactionReceipt(ctx, common.HexToHash(hash))

	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}

	log.Println(transactionReceipt)
}
