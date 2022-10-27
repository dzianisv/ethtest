package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthclientApiProvider struct {
	client   *ethclient.Client
	endpoint string
}

func NewEthcelintApiProvider(endpoint string) (*EthclientApiProvider, error) {
	rpcClient, err := rpc.DialHTTP(endpoint)

	if err != nil {
		return nil, err
	}

	rpcClient.SetHeader("User-Agent", "ethtest")

	client := ethclient.NewClient(rpcClient)

	return &EthclientApiProvider{
		client:   client,
		endpoint: endpoint,
	}, nil
}

func (self *EthclientApiProvider) eth_transactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return self.client.TransactionReceipt(ctx, hash)
}

func (self *EthclientApiProvider) web3_clientVersion(ctx context.Context) (string, error) {
	return "", fmt.Errorf("Not supported")
}
