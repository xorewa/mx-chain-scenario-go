package worldmock

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/stretchr/testify/require"
)

func TestUpdateBalanceWithDeltaRejectsNegativeResultWithoutMutation(t *testing.T) {
	world := NewMockWorld()
	address := bytes.Repeat([]byte{0x01}, len(core.SystemAccountAddress))
	account := world.AcctMap.CreateAccount(address, world)
	account.Balance = big.NewInt(10)

	err := world.UpdateBalanceWithDelta(address, big.NewInt(-11))
	require.Error(t, err)
	require.Equal(t, int64(10), account.Balance.Int64())

	require.NoError(t, world.UpdateBalanceWithDelta(address, big.NewInt(-10)))
	require.Zero(t, account.Balance.Sign())
	require.NoError(t, world.UpdateBalanceWithDelta(address, big.NewInt(5)))
	require.Equal(t, int64(5), account.Balance.Int64())
}

func TestGenerateMockAddressCheckedRejectsMalformedInputs(t *testing.T) {
	validAddress := bytes.Repeat([]byte{0x01}, len(core.SystemAccountAddress))

	_, err := GenerateMockAddressChecked(validAddress[:len(validAddress)-1], 1, []byte{0x05, 0x00})
	require.Error(t, err)
	_, err = GenerateMockAddressChecked(validAddress, 1, nil)
	require.Error(t, err)
	_, err = GenerateMockAddressChecked(validAddress, 1, []byte{0x05})
	require.Error(t, err)

	generated, err := GenerateMockAddressChecked(validAddress, 1, []byte{0x05, 0x00})
	require.NoError(t, err)
	require.Len(t, generated, len(core.SystemAccountAddress))
}

func TestMockWorldNewAddressReturnsMalformedInputError(t *testing.T) {
	world := NewMockWorld()
	_, err := world.NewAddress([]byte{0x01}, 1, []byte{0x05, 0x00})
	require.Error(t, err)
}
