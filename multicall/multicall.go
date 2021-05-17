package multicall

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type Multicall interface {
	CallRaw(calls ViewCalls, block string) (*Result, error)
	Call(calls ViewCalls, block string) (*Result, error)
	Contract() string
}

type multicall struct {
	eth    ETHIfc
	config *Config
}

type ETHIfc interface {
	MakeEthRpcCall(cntAddress, data string, blockNumber int) (string, error)
}

func New(eth ETHIfc, opts ...Option) (Multicall, error) {
	config := &Config{
		MulticallAddress: MainnetAddress,
		Gas:              "0x400000000",
	}

	for _, opt := range opts {
		opt(config)
	}

	return &multicall{
		eth:    eth,
		config: config,
	}, nil
}

type CallResult struct {
	Success bool
	Raw     []byte
	Decoded []interface{}
}

type Result struct {
	BlockNumber uint64
	Calls       map[string]CallResult
}

const AggregateMethod = "0x17352e13"

func (mc multicall) CallRaw(calls ViewCalls, block string) (*Result, error) {
	resultRaw, err := mc.makeRequest(calls, block)
	if err != nil {
		return nil, err
	}
	return calls.decodeRaw(resultRaw)
}

func (mc multicall) Call(calls ViewCalls, block string) (*Result, error) {
	resultRaw, err := mc.makeRequest(calls, block)
	if err != nil {
		return nil, err
	}
	return calls.decode(resultRaw)
}

func (mc multicall) makeRequest(calls ViewCalls, block string) (string, error) {
	payloadArgs, err := calls.callData()
	if err != nil {
		return "", err
	}
	data := AggregateMethod + hex.EncodeToString(payloadArgs)
	blockNo, err := strconv.ParseInt(strings.TrimPrefix(block, "0x"), 16, 64)
	if err != nil {
		return "", fmt.Errorf("invalid block no: %s", block)
	}
	return mc.eth.MakeEthRpcCall(mc.config.MulticallAddress, data, int(blockNo))
}

func (mc multicall) Contract() string {
	return mc.config.MulticallAddress
}
