# Multix

A Go library for interacting with [Multicall3](https://github.com/mds1/multicall) contracts on Ethereum. This library provides a clean and efficient way to batch multiple contract calls into a single transaction.

## Features

- **Batch Contract Calls**: Execute multiple contract calls in a single transaction
- **Type Safety**: Full Go type safety with ABI integration
- **Easy to Use**: Simple API for adding and executing calls
- **Ethereum Compatible**: Built on top of go-ethereum
- **Gas Efficient**: Reduces gas costs by batching calls

## Installation

```bash
go get github.com/0xjakena4/multix
```

## Quick Start

```go
package main

import (
    "log"
    "os"
    
    multical3 "github.com/0xjakena4/multix"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // Connect to Ethereum network
    client, err := ethclient.Dial(os.Getenv("RPC_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Multicall3 contract address (mainnet)
    multicall3Address := common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")
    
    // Create multicall instance
    mc := multical3.New(multicall3Address, client)
    caller := mc.NewCaller()

    // Add your calls here...
    // See examples below
}
```

## Usage

### Basic Example

The following example demonstrates how to batch multiple ERC20 token calls:

```go
// Define your contract ABI
const tokenABI = `[
    {"constant": true, "inputs": [], "name": "symbol", "outputs": [{"internalType": "string", "name": "", "type": "string"}], "payable": false, "stateMutability": "view", "type": "function"},
    {"constant": true, "inputs": [], "name": "decimals", "outputs": [{"internalType": "uint8", "name": "", "type": "uint8"}], "payable": false, "stateMutability": "view", "type": "function"},
    {"constant": true, "inputs": [{"internalType": "address", "name": "", "type": "address"}], "name": "balanceOf", "outputs": [{"internalType": "uint256", "name": "", "type": "uint256"}], "payable": false, "stateMutability": "view", "type": "function"}
]`

// Parse ABI
tokenABIParsed, err := abi.JSON(strings.NewReader(tokenABI))
if err != nil {
    log.Fatal(err)
}

// Token and user addresses
tokenAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48") // USDC
userAddress := common.HexToAddress("0xd166B8Fca31A21962AAE30EBaaA2B5464e8Ea2B3")

// Create multicall instance
mc := multical3.New(multicall3Address, client)
caller := mc.NewCaller()

// Prepare output variables
var symbol string
var decimals uint8
var balance *big.Int

// Add calls to the batch
caller.AddByABI(tokenAddress, &tokenABIParsed, "symbol", &symbol)
caller.AddByABI(tokenAddress, &tokenABIParsed, "decimals", &decimals)
caller.AddByABI(tokenAddress, &tokenABIParsed, "balanceOf", &balance, userAddress)

// Execute the batch
err = caller.Aggregate()
if err != nil {
    log.Fatal(err)
}

// Use the results
fmt.Printf("Token Symbol: %s\n", symbol)
fmt.Printf("Token Decimals: %d\n", decimals)
fmt.Printf("User Balance: %s\n", balance.String())
```

## API Reference

### Multical3

The main struct for interacting with Multicall3 contracts.

```go
type Multical3 struct {
    Contract
}
```

#### New(address, client)

Creates a new Multical3 instance.

```go
func New(address common.Address, client *ethclient.Client) *Multical3
```

#### NewCaller()

Creates a new Caller for batching calls.

```go
func (m *Multical3) NewCaller() *Caller
```

### Caller

The Caller struct manages batched calls.

```go
type Caller struct {
    multical3  *Multical3
    targets    []common.Address
    funcPack   []FuncPack
    funcUnpack []FuncUnpack
    in         [][]any
    out        []any
}
```

#### AddByABI(target, abi, function, output, inputs...)

Adds a call to the batch.

```go
func (c *Caller) AddByABI(target common.Address, abi *abi.ABI, f string, out any, in ...any) *Caller
```

Parameters:
- `target`: The contract address to call
- `abi`: The contract ABI
- `f`: The function name to call
- `out`: Pointer to the output variable
- `in`: Function input parameters

#### Aggregate()

Executes all batched calls in a single transaction.

```go
func (c *Caller) Aggregate() error
```

## Examples

See the `example/` directory for complete working examples:

- `ex1.go`: Basic ERC20 token calls (symbol, decimals, balanceOf)

## Requirements

- Go 1.24.3 or later
- go-ethereum v1.15.11

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 