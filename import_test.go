package porto

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddImportPathAddsVanityImport(t *testing.T) {
	cwd, _ := os.Getwd()
	hasChanged, newContent, err := addImportPath(
		cwd+"/testdata/leftpad/leftpad.go",
		"mypackage")

	require.NoError(t, err)
	assert.True(t, hasChanged)
	assert.Equal(t, "package leftpad // import \"mypackage\"", string(newContent[14:51]))
}

func TestAddImportPathFixesTheVanityImport(t *testing.T) {
	cwd, _ := os.Getwd()
	hasChanged, newContent, err := addImportPath(
		cwd+"/testdata/rightpad/rightpad.go",
		"mypackage")

	require.NoError(t, err)
	assert.True(t, hasChanged)
	assert.Equal(t, "package rightpad // import \"mypackage\"", string(newContent[:38]))
}

func TestIsIgnoredFile(t *testing.T) {
	assert.True(
		t,
		isIgnoredFile(
			[]*regexp.Regexp{regexp.MustCompile(".*\\.pb\\.go$")},
			"myfile.pb.go",
		),
	)

	assert.False(
		t,
		isIgnoredFile(
			[]*regexp.Regexp{},
			"myfile.pb.go",
		),
	)
}
