package modedit

import (
	_ "embed"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/expected_output.mod
var expectedOutput []byte

func TestCopyReplace(t *testing.T) {
	out := filepath.Join(t.TempDir(), "new.mod")
	err := CopyReplace("testdata/insights.mod", "testdata/sandbox.mod", out)
	require.NoError(t, err)

	res, err := ioutil.ReadFile(out)
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, res)
}
