package scenfileresolver

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveAbsolutePathWithinContext(t *testing.T) {
	contextFile := filepath.Join(t.TempDir(), "scenarios", "root.scen.json")
	resolver := NewDefaultFileResolver().WithContext(contextFile)

	fullPath, err := resolver.ResolveAbsolutePath("nested/test.wasm")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(filepath.Dir(contextFile), "nested", "test.wasm"), fullPath)
}

func TestResolveAbsolutePathRejectsTraversalOutsideContext(t *testing.T) {
	contextFile := filepath.Join(t.TempDir(), "scenarios", "root.scen.json")
	resolver := NewDefaultFileResolver().WithContext(contextFile)

	_, err := resolver.ResolveAbsolutePath("../../etc/passwd")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrPathEscapesContext))
}

func TestResolveAbsolutePathAllowsArtifactWithinRustProject(t *testing.T) {
	projectDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(projectDir, "Cargo.toml"), []byte("[package]\nname = \"contract\"\n"), 0o600))
	contextFile := filepath.Join(projectDir, "scenarios", "root.scen.json")
	require.NoError(t, os.MkdirAll(filepath.Dir(contextFile), 0o755))

	resolver := NewDefaultFileResolver().WithContext(contextFile)
	fullPath, err := resolver.ResolveAbsolutePath("../output/contract.mxsc.json")

	require.NoError(t, err)
	require.Equal(t, filepath.Join(projectDir, "output", "contract.mxsc.json"), fullPath)
}

func TestResolveAbsolutePathRejectsTraversalOutsideRustProject(t *testing.T) {
	projectDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(projectDir, "Cargo.toml"), []byte("[package]\nname = \"contract\"\n"), 0o600))
	contextFile := filepath.Join(projectDir, "scenarios", "root.scen.json")
	require.NoError(t, os.MkdirAll(filepath.Dir(contextFile), 0o755))

	resolver := NewDefaultFileResolver().WithContext(contextFile)
	_, err := resolver.ResolveAbsolutePath("../../outside-project.mxsc.json")

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrPathEscapesContext))
}

func TestResolveAbsolutePathAllowsExplicitReplacement(t *testing.T) {
	resolver := NewDefaultFileResolver().ReplacePath("contract.wasm", "../shared/contract.wasm")

	fullPath, err := resolver.ResolveAbsolutePath("contract.wasm")
	require.NoError(t, err)
	require.Equal(t, filepath.Clean("../shared/contract.wasm"), fullPath)
}

func TestResolveFileValueRejectsTraversalOutsideContext(t *testing.T) {
	contextFile := filepath.Join(t.TempDir(), "scenarios", "root.scen.json")
	resolver := NewDefaultFileResolver().WithContext(contextFile)

	_, err := resolver.ResolveFileValue("../../etc/passwd")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrPathEscapesContext))
}

func TestResolveFileValueReadsFileWithinContext(t *testing.T) {
	rootDir := t.TempDir()
	contextFile := filepath.Join(rootDir, "scenarios", "root.scen.json")
	targetDir := filepath.Dir(contextFile)
	require.NoError(t, os.MkdirAll(targetDir, 0o755))
	targetFile := filepath.Join(targetDir, "contract.wasm")
	require.NoError(t, os.WriteFile(targetFile, []byte("wasm"), 0o600))

	resolver := NewDefaultFileResolver().WithContext(contextFile)

	contents, err := resolver.ResolveFileValue("contract.wasm")
	require.NoError(t, err)
	require.Equal(t, []byte("wasm"), contents)
}
