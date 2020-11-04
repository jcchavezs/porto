package porto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindGoModule(t *testing.T) {
	module, found := findGoModule("./testdata/withgomod")
	assert.True(t, found)
	assert.Equal(t, "github.com/jcchavezs/porto/testmodule", module)

	_, found = findGoModule("./testdata/withoutgomod")
	assert.False(t, found)
}

func TestIsGoFile(t *testing.T) {
	assert.True(t, isGoFile("example.go"))
	assert.False(t, isGoFile(".go"))
}

func TestIsGoTestFile(t *testing.T) {
	assert.True(t, isGoFile("example_test.go"))
}
