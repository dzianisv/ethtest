package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Result struct {
	Id      int64         `json:"id,omitempty"`
	Jsonrpc string        `json:"jsonrpc,omitempty"`
	Result  types.Receipt `json:"result" `
}

func goTransactionReceipt(client *http.Client, endpoint string, hash common.Hash, context context.Context) (*types.Receipt, error) {
	body := []byte(fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionReceipt\",\"params\":[\"%s\"],\"id\":0}", hash.String()))

	req, err := http.NewRequestWithContext(context, "POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %s", err)
		return nil, err
	}

	defer res.Body.Close()

	result := Result{}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response body: %s", err)
		return nil, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Printf("Failed to parse body \"%s\": %s", body, err)
		return nil, err
	}
	return &result.Result, nil
}

func goClientVersion(client *http.Client, endpoint string, context context.Context) (string, error) {
	body := []byte(fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[],\"id\":0}"))

	req, err := http.NewRequestWithContext(context, "POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ethtest")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %s", err)
		return "", err
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response body: %s", err)
		return "", err
	}

	result := make(map[string]interface{})

	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Printf("Failed to parse body \"%s\": %s", body, err)
		return "", err
	}
	return result["result"].(string), nil
}
