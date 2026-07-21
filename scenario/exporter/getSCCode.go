package exporter

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetSCCode retrieves the bytecode of a WASM module from a local file.
// Callers receive file-system failures as errors so programmatic exporters do
// not terminate their host process for malformed or missing fixture paths.
func GetSCCode(fileName string) ([]byte, error) {
	code, err := os.ReadFile(filepath.Clean(fileName))
	if err != nil {
		return nil, fmt.Errorf("get SC code from %q: %w", fileName, err)
	}

	return code, nil
}
