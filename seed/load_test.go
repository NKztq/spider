

// load_test.go - UT for load.go.

package seed

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_NormalCase(t *testing.T) {
	seedPath := "./testdata/seed.json"

	expectSeeds := []string{"http://www.baidu.com", "http://www.sina.com.cn"}

	seeds, err := Load(seedPath)

	assert.NoError(t, err)
	assert.Equal(t, expectSeeds, seeds)
}

func TestLoad_InvalidFormat(t *testing.T) {
	seedPath := "./testdata/seed1.json"

	_, err := Load(seedPath)

	assert.True(t, strings.Contains(err.Error(), "cannot unmarshal string into Go value of type []string"))
}

func TestLoad_DropDuplicate(t *testing.T) {
	seedPath := "./testdata/seed.json"

	seeds, err := Load(seedPath)
	assert.NoError(t, err)
	assert.Len(t, seeds, 2)
}

func TestLoad_FileNotExist(t *testing.T) {
	// not exist
	seedPath := "./testdata/seed100.json"

	_, err := Load(seedPath)

	assert.True(t, strings.Contains(err.Error(), "no such file or directory"))
}
