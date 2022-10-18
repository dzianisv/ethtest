package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func goTransactionReceipt(endpoint string, hash common.Hash, context context.Context) (*types.Receipt, error) {
	body := []byte(fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionreceipt\",\"params\":[\"%s\"],\"id\":0}", hash.String()))
	req, err := http.NewRequestWithContext(context, "POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	receipt := types.Receipt{}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, receipt)
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}
