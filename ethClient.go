package main

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func ethTransactionReceipt(endpoint string, hash common.Hash, ctx context.Context) (*types.Receipt, error) {

	rpcClient, err := rpc.DialHTTP(endpoint)

	if err != nil {
		return nil, err
	}

	// userAgent := fmt.Sprintf("ethtest %d %s", i, time.Now().Format(time.RFC3339))
	// rpcClient.SetHeader("User-Agent", userAgent)

	client := ethclient.NewClient(rpcClient)
	return client.TransactionReceipt(ctx, hash)
}
