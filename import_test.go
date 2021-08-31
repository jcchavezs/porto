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

func TestFindFilesWithVanityImport(t *testing.T) {
	cwd, _ := os.Getwd()

	t.Run("one file listed", func(t *testing.T) {
		c, err := findAndAddVanityImportForModuleDir(
			cwd,
			cwd+"/testdata/leftpad",
			"github.com/jcchavezs/porto-integration-leftpad",
			Options{
				ListDiffFiles: true,
			},
		)

		require.NoError(t, err)
		assert.Equal(t, 2, c)
	})

	t.Run("no files listed", func(t *testing.T) {
		c, err := findAndAddVanityImportForModuleDir(
			cwd,
			cwd+"/testdata/nopad",
			"github.com/jcchavezs/porto-integration/nopad",
			Options{
				ListDiffFiles: true,
			},
		)

		require.NoError(t, err)
		assert.Equal(t, 0, c)
	})

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
