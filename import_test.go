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
	assert.Equal(t, "package leftpad // import \"mypackage\"", string(newContent[15:52]))
}

func TestAddImportAutogenerated(t *testing.T) {
	cwd, _ := os.Getwd()
	hasChanged, _, err := addImportPath(
		cwd+"/testdata/codegen/generated.go",
		"codegen")

	assert.Equal(t, errGenerated, err)
	assert.False(t, hasChanged)
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
			cwd+"/testdata/nopad",
			"github.com/jcchavezs/porto-integration/nopad",
			Options{
				ListDiffFiles: true,
			},
		)

		require.NoError(t, err)
		assert.Equal(t, 0, c)
	})

	t.Run("skip file", func(t *testing.T) {
		c, err := findAndAddVanityImportForModuleDir(
			cwd,
			cwd+"/testdata/leftpad",
			cwd+"/testdata/leftpad",
			"github.com/jcchavezs/porto-integration-leftpad",
			Options{
				ListDiffFiles:    true,
				SkipFilesRegexes: []*regexp.Regexp{regexp.MustCompile(`leftpad\.go`)},
			},
		)

		require.NoError(t, err)
		assert.Equal(t, 1, c)
	})

	t.Run("skip dir", func(t *testing.T) {
		c, err := findAndAddVanityImportForModuleDir(
			cwd,
			cwd+"/testdata",
			cwd+"/testdata",
			"github.com/jcchavezs/porto/integration",
			Options{
				ListDiffFiles: true,
				SkipDirsRegexes: []*regexp.Regexp{
					regexp.MustCompile(`^codegen$`),
					regexp.MustCompile(`^leftpad$`),
					regexp.MustCompile(`^rightpad$`),
				},
			},
		)

		require.NoError(t, err)
		assert.Equal(t, 2, c)
	})

	t.Run("restrict to files", func(t *testing.T) {
		c, err := findAndAddVanityImportForModuleDir(
			cwd,
			cwd+"/testdata/leftpad",
			cwd+"/testdata/leftpad",
			"github.com/jcchavezs/porto-integration-leftpad",
			Options{
				ListDiffFiles:          true,
				RestrictToFilesRegexes: []*regexp.Regexp{regexp.MustCompile(`^other\.go$`)},
			},
		)

		require.NoError(t, err)
		assert.Equal(t, 1, c)
	})

	t.Run("restrict to dir", func(t *testing.T) {
		c, err := findAndAddVanityImportForModuleDir(
			cwd,
			cwd+"/testdata",
			cwd+"/testdata",
			"github.com/jcchavezs/porto/integration",
			Options{
				ListDiffFiles:         true,
				RestrictToDirsRegexes: []*regexp.Regexp{regexp.MustCompile(`^withoutgomod`)},
			},
		)

		require.NoError(t, err)
		assert.Equal(t, 2, c)
	})

	t.Run("skip and include file", func(t *testing.T) {
		c, err := findAndAddVanityImportForModuleDir(
			cwd,
			cwd+"/testdata/leftpad",
			cwd+"/testdata/leftpad",
			"github.com/jcchavezs/porto-integration-leftpad",
			Options{
				ListDiffFiles:          true,
				RestrictToFilesRegexes: []*regexp.Regexp{regexp.MustCompile(`other\.go`)},
				SkipFilesRegexes:       []*regexp.Regexp{regexp.MustCompile(`leftpad\.go`)},
			},
		)

		require.NoError(t, err)
		assert.Equal(t, 1, c)
	})

}

func TestMatchesAny(t *testing.T) {
	assert.True(
		t,
		matchesAny(
			[]*regexp.Regexp{regexp.MustCompile(".*\\.pb\\.go$")},
			"myfile.pb.go",
		),
	)

	assert.False(
		t,
		matchesAny(
			[]*regexp.Regexp{},
			"myfile.pb.go",
		),
	)

	assert.True(
		t,
		matchesAny(
			[]*regexp.Regexp{regexp.MustCompile("^third_party$")},
			"third_party",
		),
	)
}

func TestIsUnexportedModule(t *testing.T) {
	assert.True(t, isUnexportedModule("go.opentelemetry.io/otel/internal", false))
	assert.True(t, isUnexportedModule("go.opentelemetry.io/otel/internal/metric", false))
	assert.False(t, isUnexportedModule("go.opentelemetry.io/otel/internalmetric", false))
	assert.False(t, isUnexportedModule("go.opentelemetry.io/otel/internal/metric", true))
}
