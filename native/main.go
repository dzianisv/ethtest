package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	timeout               = 45
	keepAlive             = 30
	maxIdleConns          = 1
	idleConnTimeout       = 60
	TLSHandshakeTimeout   = 10
	expectContinueTimeout = 1
)

type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type JSONRPCErrorResponse struct {
	JSONRPC string       `json:"jsonrpc"`
	ID      int          `json:"id"`
	Error   JSONRPCError `json:"error"`
}

func (j JSONRPCError) Error() string {
	return fmt.Sprintf(
		`JSONRPCError code=%d message="%s"`,
		j.Code,
		j.Message,
	)
}

func main() {
	var (
		transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   timeout * time.Second,
				KeepAlive: keepAlive * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     false,
			MaxIdleConns:          maxIdleConns,
			IdleConnTimeout:       idleConnTimeout * time.Second,
			TLSHandshakeTimeout:   TLSHandshakeTimeout * time.Second,
			ExpectContinueTimeout: expectContinueTimeout * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		}

		transport2 = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   timeout * time.Second,
				KeepAlive: keepAlive * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     false,
			MaxIdleConns:          maxIdleConns,
			IdleConnTimeout:       idleConnTimeout * time.Second,
			TLSHandshakeTimeout:   TLSHandshakeTimeout * time.Second,
			ExpectContinueTimeout: expectContinueTimeout * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		}
		transport3 = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   timeout * time.Second,
				KeepAlive: keepAlive * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     false,
			MaxIdleConns:          maxIdleConns,
			IdleConnTimeout:       idleConnTimeout * time.Second,
			TLSHandshakeTimeout:   TLSHandshakeTimeout * time.Second,
			ExpectContinueTimeout: expectContinueTimeout * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		}
	)

	clients := []*http.Client{
		&http.Client{
			Timeout:   timeout * time.Second,
			Transport: transport,
		},
		&http.Client{
			Timeout:   timeout * time.Second,
			Transport: transport2,
		},
		&http.Client{
			Timeout:   timeout * time.Second,
			Transport: transport3,
		},
	}

	errors_n := 0
	requests_n := 0
	for {
		for i, client := range clients {
			requests_n += 1
			res, err := CallJSONRPCMethod(client)
			if err != nil {
				errors_n += 1
				log.Printf("Client %d: %s: %s", i, res, err)
			}
		}

		rand.Seed(time.Now().UnixNano())
		wait_time := 1 + rand.Intn(130)
		log.Printf("Requests: %d, errors: %d, random wait: %ds", requests_n, errors_n, wait_time)
		time.Sleep(time.Duration(wait_time) * time.Second)
	}
}

func CallJSONRPCMethod(c *http.Client) (res interface{}, err error) {
	jsonBody, err := json.Marshal(JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		ID:      1,
	})

	if err != nil {
		return nil, err
	}

	host := os.Getenv("HOST")
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	if len(host) == 0 || len(name) == 0 || len(password) == 0 {
		log.Fatalf("Host, name, password are not set")
	}

	authorization := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", name, password)))

	url := fmt.Sprintf("https://%s", host)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go/net/http/test")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authorization))

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	rpcResp := &JSONRPCResponse{}
	rpcErrResp := &JSONRPCResponse{}
	err = json.NewDecoder(resp.Body).Decode(rpcResp)
	if err != nil {
		err = json.NewDecoder(resp.Body).Decode(rpcErrResp)
		if err != nil {
			return nil, err
		}
		return nil, rpcErrResp.Error
	}

	return rpcResp.Result, nil
}
