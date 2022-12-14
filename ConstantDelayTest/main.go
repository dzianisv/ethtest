package main

import (
	"bytes"
	"encoding/json"
	"flag"
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
	keepAlive             = 30
	maxIdleConns          = 1000
	idleConnTimeout       = 60
	TLSHandshakeTimeout   = 10
	expectContinueTimeout = 1
)

func main() {
	timeoutFlag := flag.Int("t", 5, "timeout")
	flag.Parse()

	timeout := *timeoutFlag

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
			KeepAlive: keepAlive * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          maxIdleConns,
		IdleConnTimeout:       idleConnTimeout * time.Second,
		TLSHandshakeTimeout:   TLSHandshakeTimeout * time.Second,
		ExpectContinueTimeout: expectContinueTimeout * time.Second,
	}

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
	errors_n := 0
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
			errors_n += 1
			fmt.Printf("Failed to query: %s\n", err)
			fmt.Printf("Errors: %d\n", errors_n)
			continue
		}

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			errors_n += 1
			log.Printf("failed to read a response: %s\n", err)
			log.Printf("erors: %d\n", errors_n)
			continue
		}

		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(body, &jsonMap)
		if err != nil {
			log.Printf("failed to decoded response\n: %s", err)
		}
	}
}
