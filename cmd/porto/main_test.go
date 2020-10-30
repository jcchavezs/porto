package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindGoMod(t *testing.T) {
	module, found := findGoMod("../../testdata/withgomod")
	assert.True(t, found)
	assert.Equal(t, "github.com/jcchavezs/porto/testmodule", module)

	_, found = findGoMod("../../testdata/withoutgomod")
	assert.False(t, found)
}

func TestIsGoFile(t *testing.T) {
	assert.True(t, isGoFile("example.go"))
	assert.False(t, isGoFile(".go"))
}

func TestIsGoTestFile(t *testing.T) {
	assert.True(t, isGoFile("example_test.go"))
}
