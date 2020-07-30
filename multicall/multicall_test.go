package multicall_test

import (
	"encoding/json"
	"fmt"
	"github.com/HAL-xyz/ethrpc"
	"github.com/HAL-xyz/web3-multicall-go/multicall"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExampleViwCall(t *testing.T) {

	cli := ethrpc.New("https://mainnet.infura.io/v3/17ed7fe26d014e5b9be7dfff5368c69d")
	mc, _ := multicall.New(cli)

	vc := multicall.NewViewCall(
		"key.1",
		"0x6b175474e89094c44da98b954eedeac495271d0f",
		"symbol()(string)",
		[]interface{}{},
	)

	vc2 := multicall.NewViewCall(
		"key.2",
		"0x6b175474e89094c44da98b954eedeac495271d0f",
		"decimals()(uint8)",
		[]interface{}{},
	)

	vcs := multicall.ViewCalls{vc, vc2}

	block := "latest"
	res, err := mc.Call(vcs, block)
	if err != nil {
		panic(err)
	}

	resJson, _ := json.Marshal(res)
	fmt.Println(string(resJson))
	fmt.Println(res)
	fmt.Println(err)

	assert.NoError(t, err)
	assert.Equal(t, res.Calls["key.1"].Decoded[0].(string), "DAI")
	assert.Equal(t, res.Calls["key.2"].Decoded[0].(uint8), uint8(18))
}
