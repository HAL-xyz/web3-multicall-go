package multicall

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestViewCall(t *testing.T) {
	vc := ViewCall{
		id:        "key",
		target:    "0x0",
		method:    "balanceOf(address, uint64)(int256)",
		arguments: []interface{}{common.HexToAddress("0x0000000000000000000000000000000000000000"), uint64(12)},
	}
	expectedArgTypes := []string{"address", "uint64"}
	expectedCallData := "295eaadf0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c"

	assert.Equal(t, expectedArgTypes, vc.argumentTypes())
	callData, err := vc.callData()
	assert.Nil(t, err)
	assert.Equal(t, expectedCallData, fmt.Sprintf("%x", callData))
}

func TestEncodeBytes32Argument(t *testing.T) {
	var bytes32Array = [32]uint8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

	vc1 := ViewCall{
		id:        "key",
		target:    "0x0",
		method:    "balanceOfPartition(bytes32, uint256)(int256)",
		arguments: []interface{}{bytes32Array, big.NewInt(12312312312313)},
	}

	_, err1 := vc1.argsCallData()
	assert.Nil(t, err1)
}
