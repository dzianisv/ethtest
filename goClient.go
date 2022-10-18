package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Result struct {
	Id      int64         `json:"id,omitempty"`
	Jsonrpc string        `json:"jsonrpc,omitempty"`
	Result  types.Receipt `json:"result" `
}

func goTransactionReceipt(endpoint string, hash common.Hash, context context.Context) (*types.Receipt, error) {
	body := []byte(fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionReceipt\",\"params\":[\"%s\"],\"id\":0}", hash.String()))

	req, err := http.NewRequestWithContext(context, "POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	var netTransport = http.DefaultTransport.(*http.Transport).Clone()
	netTransport.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}
	netTransport.ForceAttemptHTTP2 = false

	client := &http.Client{
		Transport: netTransport,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	result := Result{}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result.Result, nil
}
