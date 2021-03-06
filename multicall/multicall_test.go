package multicall_test

import (
	"fmt"
	"github.com/HAL-xyz/ethrpc"
	"github.com/HAL-xyz/web3-multicall-go/multicall"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var cli = GetCustomCli()
var lastBlockNo, _ = cli.cli.EthBlockNumber()
var lastBlockHex = fmt.Sprintf("%x", lastBlockNo)

func GetCustomCli() CustomCli {
	return CustomCli{cli: ethrpc.New("https://mainnet.infura.io/v3/17ed7fe26d014e5b9be7dfff5368c69d")}
}

type CustomCli struct {
	cli *ethrpc.EthRPC
}

func (c CustomCli) MakeEthRpcCall(cntAddress, data string, blockNumber int) (string, error) {
	params := ethrpc.T{
		To:   cntAddress,
		From: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		Data: data,
	}
	hexBlockNo := fmt.Sprintf("0x%x", blockNumber)

	return c.cli.EthCall(params, hexBlockNo)
}

func TestExampleViwCall(t *testing.T) {

	vc := multicall.NewViewCall(
		"key.1",
		"0x5d3a536E4D6DbD6114cc1Ead35777bAB948E3643",
		"totalReserves()(uint256)",
		[]interface{}{},
	)
	vcs := multicall.ViewCalls{vc}
	mc, _ := multicall.New(cli)
	res, err := mc.Call(vcs, lastBlockHex)
	assert.NoError(t, err)

	assert.IsType(t, big.NewInt(0), res.Calls["key.1"].Decoded[0])
}

func TestViwCallWithDecodeError(t *testing.T) {
	vc1 := multicall.NewViewCall(
		"vc1-ok",
		"0x1f9840a85d5af5bf1d1762f925bdaddc4201f984",
		"symbol()(string)",
		[]interface{}{},
	)
	vc2 := multicall.NewViewCall(
		"vc2-ok",
		"0xa417221ef64b1549575c977764e651c9fab50141",
		"latestAnswer()(int256)",
		[]interface{}{},
	)
	vc3 := multicall.NewViewCall(
		"vc3-fail",
		"0x7a250d5630b4cf539739df2c5dacb4c659f2488d",
		"getAmountsOut(uint256,address[])(uint256[])",
		[]interface{}{big.NewInt(1000000000), []common.Address{common.HexToAddress("0x9ceb84f92a0561fa3cc4132ab9c0b76a59787544"), common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")}},
	)
	vc4 := multicall.NewViewCall(
		"vc4-ok",
		"0x7a250d5630b4cf539739df2c5dacb4c659f2488d",
		"getAmountsOut(uint256,address[])(uint256[])",
		[]interface{}{big.NewInt(1000000000), []common.Address{common.HexToAddress("0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"), common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")}},
	)
	vc5 := multicall.NewViewCall(
		"vc5-ok",
		"0x09cabec1ead1c0ba254b09efb3ee13841712be14",
		"getTokenToEthOutputPrice(uint256)(uint256)",
		[]interface{}{big.NewInt(1)},
	)

	vcs := multicall.ViewCalls{vc1, vc2, vc3, vc4, vc5}
	mc, _ := multicall.New(GetCustomCli())
	res, err := mc.Call(vcs, lastBlockHex)
	assert.NoError(t, err)

	assert.Equal(t, true, res.Calls["vc1-ok"].Success)
	assert.Equal(t, "UNI", res.Calls["vc1-ok"].Decoded[0].(string))
	assert.Equal(t, false, res.Calls["vc2-fail"].Success)
	assert.Equal(t, false, res.Calls["vc3-fail"].Success)
	assert.Equal(t, true, res.Calls["vc4-ok"].Success)
	assert.Equal(t, true, res.Calls["vc5-ok"].Success)
}
