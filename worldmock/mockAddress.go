package worldmock

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// GenerateMockAddressChecked simulates creation of a new address by the
// protocol and validates the protocol-shaped inputs before indexing them.
func GenerateMockAddressChecked(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	if len(creatorAddress) != len(core.SystemAccountAddress) {
		return nil, fmt.Errorf("invalid creator address length: got %d, expected %d", len(creatorAddress), len(core.SystemAccountAddress))
	}
	if len(vmType) != core.VMTypeLen {
		return nil, fmt.Errorf("invalid VM type length: got %d, expected %d", len(vmType), core.VMTypeLen)
	}

	result := make([]byte, len(core.SystemAccountAddress))
	result[10] = 0x11
	result[11] = 0x11
	result[12] = 0x11
	result[13] = 0x11
	copy(result[14:29], creatorAddress)
	result[29] = byte(creatorNonce)
	copy(result[30:], creatorAddress[30:])
	copy(result[vmcommon.NumInitCharactersForScAddress-core.VMTypeLen:], vmType)
	return result, nil
}

// GenerateMockAddress preserves the historical panic-based helper for
// external compatibility. Internal scenario execution uses the checked API.
//
// Not an actual blockchain hook, just a helper method.
func GenerateMockAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) []byte {
	result, err := GenerateMockAddressChecked(creatorAddress, creatorNonce, vmType)
	if err != nil {
		panic(fmt.Sprintf("GenerateMockAddress: %v", err))
	}
	return result
}
