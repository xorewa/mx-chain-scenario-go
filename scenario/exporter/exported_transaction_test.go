package exporter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDeployTransactionRejectsInvalidCodePathWithoutPanic(t *testing.T) {
	testCases := []string{
		"",
		"wasm/module.wasm",
		"file:",
		"file:this-file-does-not-exist.wasm",
	}

	for _, codePath := range testCases {
		t.Run(codePath, func(t *testing.T) {
			tx, err := CreateDeployTransaction(nil, codePath, []byte("sender"), 1, 1)
			require.Nil(t, tx)
			require.Error(t, err)
		})
	}
}

func TestGetSCCodeReturnsReadError(t *testing.T) {
	code, err := GetSCCode("this-file-does-not-exist.wasm")
	require.Nil(t, code)
	require.Error(t, err)
}
