package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ApiProvider interface {
	web3_clientVersion(context.Context) (string, error)
	eth_transactionReceipt(context.Context, common.Hash) (*types.Receipt, error)
}

type Response struct {
	I         int
	Err       error
	DelayMili int64
}

var transactions = []string{
	"0xf5de9f2bd29a1db6c4e5c4e0fd353c464ee47a890fdc38379da2141792a16214",
	"0x26901786c7ed5009a6785aef80e67f2e7494757d27065bcb68cbd4dd41eb42f6",
	"0x3c086163dfc946c89a64d22499c3688b4ce46575285d7e9693b683f0395af04c",
	"0xe90e15046653af7ea20fbb2cdab196956256d14514e9e56e4d945960ddc57a55",
	"0xa0cd5f7428d2ec9758ba30b93cfcc2d7f274d09e8fb0196b77eac3ed1c0138bf",
	"0xf9cee393b68c4b2df3966fed6adc0a50048a6a81a1274294cf78d81d874025a8",
	"0xf2288ef04c1486fca0072b22c5b8506d949ca608cb30888b9c90eb3abbf55082",
	"0x081f731b9f597e4a46dd5b8745fc8eae4da33ae0e2aaf9d5214e7f3619e2c149",
	"0x5901c20b9ff1941200b4246c4bafaa27ce52812437bb0cf25dbdba1894d628dd",
	"0xe6041428d502d1ba9dade416df508b70956f70478d29a1c9e969c8e0441e6452",
	"0xcc401e728c7d8209bc4881a584d5860007116e9c8915f819afbb895784b49a8f",
	"0x3de5a6c946e657f02f705025cc950865d48d6001c6f546fe75304ae323df4ee8",
	"0x3b644fd621d1f3a568382cee3fe0bca36fa70849cbc71abae0afa4cad98524fc",
	"0x613bb66c69be561188af08344da7f3b9dc5df8f735e05b82ae79681aabf286a1",
	"0x9faed3f7652d42918cedffdafdac9afae1ff0c9cf7712007605b9b244623bc0c",
	"0xc1830eac12a6c9f3dfda8ff0768b24df47c42715fcc312dcaf0e2a1b9fb92be6",
	"0xcf14f6a34f53a3ba4efc7dee2f6f3191a67e421a33c3ececc5c2c23f0531b1d1",
	"0x78221362ac0648cef2432abd52e0eb44ff07c9fb147beb1eb3af6c4d272da464",
	"0xce231892a48622f3afe031643e5f115f4248dd22388de5836b204a3620dde496",
	"0xd7c0b072b3fab99cac072d1476ed9561840362badc7e75e3283003f87d24571a",
	"0x16b60490637d6f08ddf1899890532910daf345efe7c953d3478ea490232e3a05",
	"0x58a9a1bf16572d870263bbb5738b6a2d67c1cdf62a6e600efc2f8b3a4779f412",
	"0xfce6b6894f87a35031ef0fe6f5dbb39822016023a45ac2dc81f7672c6c6101c8",
	"0x00fc0912877a8c25ba21d8e77bcec8c5461e4223b7b6ab846958bb35ba375255",
	"0xa3fa06c59212cbee9cb9629a7a333eb4e21a53795dd153fafaf4a8f08d62e274",
	"0x956b12178efcfbc066eb2693ce32d3ab35b855ce0bf45e8ac7d9a5aa470d4cd5",
	"0x4acefb98484a46cf015f69338db1d14cc23edcb1d21f2f3107d4ff3cfb14b1d1",
	"0xc4a56bed5abf8103715dd616e960947a1aca7e8fa05dcc2d4275fcedba712eda",
	"0x318ad867f14cbebff79032b4555dfc4c14631c9a586d041b28455ae97e245227",
	"0x14394cbedfd05d8ddb96c0fce920ccc1fa656c6eb00416b002871430c6e4c7ce",
}

func query(ApiProvider ApiProvider, url string, i int, c chan Response, timeout int, apiMethod string) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(timeout)*time.Second))
	defer cancel()

	start_t := time.Now()
	var err error = nil

	if apiMethod == "web3_clientVersion" {
		_, err = ApiProvider.web3_clientVersion(ctx)
	} else if apiMethod == "eth_transactionReceipt" {
		tx := transactions[i%len(transactions)]
		_, err = ApiProvider.eth_transactionReceipt(ctx, common.HexToHash(tx))
	}

	delay_ms := time.Now().UnixMilli() - start_t.UnixMilli()

	if err != nil {
		c <- Response{i, err, delay_ms}
		return
	}

	c <- Response{i, nil, delay_ms}
}

func ethclientTest() {

}

func queryEthclient(client *ethclient.Client, i int, feedback chan Response, apiMethod string) {

}

func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

var validMethods = []string{"web3_clientVersion", "eth_transactionReceipt"}

func isSupportedMethod(method string) bool {
	for _, m := range validMethods {
		if m == method {
			return true
		}
	}

	return false
}

type BreakdownEntry struct {
	lessThan    int64
	count       int
	percenntage float32
}

func percentileBreakdown(data []int64) ([]BreakdownEntry, int64, int64) {
	max_entry := data[0]
	sum_ms := int64(0)

	breakdowns := make([]BreakdownEntry, 10)

	for i, entry := range data {
		max_entry = Max(max_entry, entry)
		sum_ms += entry
		log.Printf("%d\t%d\n", i, max_entry)
	}

	avg_ms := sum_ms / int64(len(data))
	divider := max_entry / int64(len(breakdowns))

	for _, item := range data {
		breakdowns[item%divider].count++
	}

	for i, breakdownEntry := range breakdowns {
		breakdownEntry.lessThan = (int64(i) + 1) * divider
		breakdownEntry.percenntage = float32(breakdownEntry.count) * 100 / float32(len(data))
	}

	return breakdowns, max_entry, avg_ms
}

func main() {
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	url := fmt.Sprintf("https://%s:%s@%s", name, password, host)

	requestFlag := flag.Int("n", 1000, "number of requests")
	concurrencyFlag := flag.Int("c", 100, "number of concurrent requetss")
	disableHttp2Flag := flag.Bool("http1", false, "disable http/2")
	apiMethodFlag := flag.String("m", "web3_clientVersion", "JSON-RPC method: <web3_clientVersion, eth_transactionReceip>")
	clientTypeFlag := flag.String("client-type", "http", "Client type: <http, ethclient>")
	timeoutFlag := flag.Int("t", 5, "request timeout")

	flag.Parse()

	if !isSupportedMethod(*apiMethodFlag) {
		log.Fatalf("Invalid JSON-RPC method: %s", *apiMethodFlag)
	}

	concurrency := *concurrencyFlag
	request_n := *requestFlag

	log.Printf("Using %s\n", url)
	if *disableHttp2Flag {
		log.Printf("Disable http/2")
	}

	c := make(chan Response)
	count := request_n
	active_n := 0
	errors_n := 0
	latency := make([]int64, request_n)

	var apiProvider ApiProvider
	var err error

	if *clientTypeFlag == "http" {
		apiProvider, err = NewHttpApiProvider(url, *disableHttp2Flag)
	} else if *clientTypeFlag == "ethclient" {
		apiProvider, err = NewEthcelintApiProvider(url)
	} else {
		log.Fatalf("Invalid Api Provier: %s", *clientTypeFlag)
	}

	if err != nil {
		log.Fatal("Failed to create an API provier: %s", err)
	}

	for {
		if active_n < concurrency && count > 0 {
			count -= 1
			active_n += 1
			go query(apiProvider, url, count, c, *timeoutFlag, *apiMethodFlag)
		} else if active_n > 0 {
			response := <-c
			latency[response.I] = response.DelayMili

			if response.Err != nil {
				errors_n += 1
				log.Printf("[%d] Error: %s", response.I, response.Err)
			}
			active_n -= 1
		} else {
			break
		}
	}

	latencyBreakdowns, max_latency, avg_latency := percentileBreakdown(latency)

	for _, breakdownEntry := range latencyBreakdowns {
		log.Printf("<%d: %d %f%%", breakdownEntry.lessThan, breakdownEntry.count, breakdownEntry.percenntage)
	}

	log.Printf("Max delay %d ms", max_latency)
	log.Printf("Average delay %d ms", avg_latency)
	log.Printf("%d/%d failed\n", errors_n, request_n)
}
