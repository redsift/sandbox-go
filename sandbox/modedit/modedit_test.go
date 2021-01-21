package modedit

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const expectedOutput=`
module github.com/redsift/sandbox-go

go 1.15

require (
	github.com/redsift/go-mangosock v0.1.2
	github.com/redsift/go-sandbox-rpc v0.1.0
	golang.org/x/mod v0.4.0
)

replace server => /run/sandbox/sift/server

replace github.com/blevesearch/bleve v1.0.14 => github.com/redsift/bleve v0.5.1-0.20201222154652-c8a1b8b0f852

replace github.com/tecbot/gorocksdb => github.com/redsift/gorocksdb v0.0.0-20180109115255-d1d69065a9b9
`

func TestCopyReplace(t *testing.T) {
	out := filepath.Join(t.TempDir(), "new.mod")
	err := CopyReplace("testdata/insights.mod", "testdata/sandbox.mod", out)
	require.NoError(t, err)

	res, err := ioutil.ReadFile(out)
	require.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(expectedOutput), strings.TrimSpace(string(res)))
}
