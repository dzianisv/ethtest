package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const GetBlockByNumberJson = `
	{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "eth_blockNumber",
		"params": []
	}`

const (
	timeout               = 30
	keepAlive             = 30
	maxIdleConns          = 1000
	idleConnTimeout       = 60
	TLSHandshakeTimeout   = 10
	expectContinueTimeout = 1
)

var transport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   timeout * time.Second,
		KeepAlive: keepAlive * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          maxIdleConns,
	IdleConnTimeout:       idleConnTimeout * time.Second,
	TLSHandshakeTimeout:   TLSHandshakeTimeout * time.Second,
	ExpectContinueTimeout: expectContinueTimeout * time.Second,
}

func main() {
	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: transport,
	}

	host := os.Getenv("HOST")
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	if len(host) == 0 {
		log.Fatalf("HOST is not set")
	}

	n := -1
	for {
		n++
		time.Sleep(1 * time.Second)
		requestBody := bytes.NewBuffer([]byte(GetBlockByNumberJson))
		req, err := http.NewRequest("POST", "https://"+host, requestBody)
		if err != nil {
			continue
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("User-Agent", fmt.Sprintf("ConstantDelayTest-%d", n))

		if len(name) != 0 {
			req.SetBasicAuth(name, password)
		}

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(n, err)
			continue
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			res.Body.Close()
			log.Printf(" %d error parsing the message: %s", n, err.Error())
			continue
		}
		res.Body.Close()

		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(body, &jsonMap)
		if err != nil {
			continue
		}
		log.Println(n, jsonMap["result"])

		continue
	}
}
