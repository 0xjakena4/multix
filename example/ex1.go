package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	multical3 "github.com/0xjakena4/multix"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const daiABI = `[
	{"constant": true, "inputs": [], "name": "symbol", "outputs": [{"internalType": "string", "name": "", "type": "string"}], "payable": false, "stateMutability": "view", "type": "function"},
	{"constant": true, "inputs": [], "name": "decimals", "outputs": [{"internalType": "uint8", "name": "", "type": "uint8"}], "payable": false, "stateMutability": "view", "type": "function"},
	{"constant": true, "inputs": [{"internalType": "address", "name": "", "type": "address"}], "name": "balanceOf", "outputs": [{"internalType": "uint256", "name": "", "type": "uint256"}], "payable": false, "stateMutability": "view", "type": "function"}
]`

var (
	multicall3Address = common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")
	usdcAddress       = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	jakeAddress       = common.HexToAddress("0xd166B8Fca31A21962AAE30EBaaA2B5464e8Ea2B3")
)

func main() {
	rpcURL := os.Getenv("MAINNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("MAINNET_RPC_URL environment variable is required")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	daiABIParsed, err := abi.JSON(strings.NewReader(daiABI))
	if err != nil {
		log.Fatalf("Failed to parse DAI ABI: %v", err)
	}

	mc := multical3.New(multicall3Address, client)
	caller := mc.NewCaller()

	var symbol string
	var decimals uint8
	var balance *big.Int

	caller.AddByABI(usdcAddress, &daiABIParsed, "symbol", &symbol)
	caller.AddByABI(usdcAddress, &daiABIParsed, "decimals", &decimals)
	caller.AddByABI(usdcAddress, &daiABIParsed, "balanceOf", &balance, jakeAddress)

	err = caller.Aggregate()
	if err != nil {
		log.Fatalf("Failed to execute multicall: %v", err)
	}

	fmt.Printf("USDC Symbol: %s\n", symbol)
	fmt.Printf("USDC Decimals: %d\n", decimals)
	fmt.Printf("Jake's %s balance: %s\n", symbol, balance.String())
}
