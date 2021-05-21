package porto

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddImportPathAddsVanityImport(t *testing.T) {
	cwd, _ := os.Getwd()
	hasChanged, newContent, err := addImportPath(
		cwd+"/testdata/leftpad/leftpad.go",
		"mypackage",
		[]string{})

	require.NoError(t, err)
	assert.True(t, hasChanged)
	assert.Equal(t, "package leftpad // import \"mypackage\"", string(newContent[14:51]))
}

func TestAddImportPathFixesTheVanityImport(t *testing.T) {
	cwd, _ := os.Getwd()
	hasChanged, newContent, err := addImportPath(
		cwd+"/testdata/rightpad/rightpad.go",
		"mypackage",
		[]string{})

	require.NoError(t, err)
	assert.True(t, hasChanged)
	assert.Equal(t, "package rightpad // import \"mypackage\"", string(newContent[:38]))
}

func TestIsToolsFile(t *testing.T) {
	isTools := isToolsFile([]byte(`
	// +build tools
	package mypackage
	`))
	assert.True(t, isTools)

	isTools = isToolsFile([]byte(`
	// +build linux,386 darwin,!cgo

	package mypackage
	`))
	assert.False(t, isTools)
}
