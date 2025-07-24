package multical3

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Contract struct {
	Address common.Address
	ABI     *abi.ABI
	Client  *ethclient.Client
}

type Multical3 struct {
	Contract
}

type FuncPack func(...any) ([]byte, error)
type FuncUnpack func(out interface{}, response []byte) error

func New(address common.Address, client *ethclient.Client) *Multical3 {
	abi, err := abi.JSON(strings.NewReader(Multicall3ABI))
	if err != nil {
		panic(err)
	}

	return &Multical3{
		Contract: Contract{
			Address: address,
			ABI:     &abi,
			Client:  client,
		},
	}
}

type Call struct {
	Target   common.Address `json:"target"`
	CallData []byte         `json:"callData"`
}

type AggregateResult struct {
	BlockNumber *big.Int
	ReturnData  [][]byte
}

type Caller struct {
	multical3  *Multical3
	targets    []common.Address
	funcPack   []FuncPack
	funcUnpack []FuncUnpack
	in         [][]any
	out        []any
}

func (m *Multical3) NewCaller() *Caller {
	return &Caller{
		multical3:  m,
		targets:    []common.Address{},
		funcPack:   []FuncPack{},
		funcUnpack: []FuncUnpack{},
		in:         [][]any{},
		out:        []any{},
	}
}

func (c *Caller) AddByABI(target common.Address, abi *abi.ABI, f string, out any, in ...any) *Caller {
	funcPack := func(args ...any) ([]byte, error) {
		return abi.Pack(f, args...)
	}
	funcUnpack := func(out interface{}, response []byte) error {
		return abi.UnpackIntoInterface(out, f, response)
	}
	c.targets = append(c.targets, target)
	c.funcPack = append(c.funcPack, funcPack)
	c.funcUnpack = append(c.funcUnpack, funcUnpack)
	c.in = append(c.in, in)
	c.out = append(c.out, out)
	return c
}

func (c *Caller) Aggregate() error {
	calls := make([]Call, len(c.funcPack))
	for i, f := range c.funcPack {
		callData, err := f(c.in[i]...)
		if err != nil {
			return err
		}
		calls[i] = Call{Target: c.targets[i], CallData: callData}
	}
	callData, err := c.multical3.ABI.Pack("aggregate", calls)
	if err != nil {
		return err
	}

	response, err := c.multical3.Client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.multical3.Address,
		Data: callData,
	}, nil)
	if err != nil {
		return err
	}

	var result AggregateResult
	err = c.multical3.ABI.UnpackIntoInterface(&result, "aggregate", response)
	if err != nil {
		return err
	}

	for i, f := range c.funcUnpack {
		err = f(c.out[i], result.ReturnData[i])
		if err != nil {
			return err
		}
	}

	return nil
}
