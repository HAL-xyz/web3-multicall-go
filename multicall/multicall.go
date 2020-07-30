package multicall

import (
	"encoding/hex"
	"github.com/HAL-xyz/ethrpc"
)

type Multicall interface {
	CallRaw(calls ViewCalls, block string) (*Result, error)
	Call(calls ViewCalls, block string) (*Result, error)
	Contract() string
}

type multicall struct {
	eth    *ethrpc.EthRPC
	config *Config
}

func New(ethCli *ethrpc.EthRPC, opts ...Option) (Multicall, error) {
	config := &Config{
		MulticallAddress: MainnetAddress,
		Gas:              "0x400000000",
	}

	for _, opt := range opts {
		opt(config)
	}

	return &multicall{
		eth:    ethCli,
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
	params := ethrpc.T{
		To: mc.config.MulticallAddress,
		From: "0x2e34c46ad2f08a66bc9ff2e9fe5918590551e958", // TODO update our lib so this isn't required
		Data: AggregateMethod + hex.EncodeToString(payloadArgs),
	}

	return mc.eth.EthCall(params, block)
}

func (mc multicall) Contract() string {
	return mc.config.MulticallAddress
}
