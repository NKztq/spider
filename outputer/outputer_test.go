// outputer_test.go - UT for outputer.go.

package outputer

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/NKztq/spider/conf"
)

func TestOutputFile(t *testing.T) {
	directory := "./test_output"
	defer func() {
		// remove UT output
		assert.NoError(t, os.RemoveAll(directory))
	}()

	fileName := "test.html"
	content := []byte("test")

	o, err := NewOutputer(conf.OutputerConf{directory, ".*.(htm|html)$"})
	assert.NoError(t, err)

	err = o.OutputFile(fileName, content)
	assert.NoError(t, err)

	fp := path.Join(directory, fileName)
	fData, err := ioutil.ReadFile(fp)
	assert.NoError(t, err)
	assert.Equal(t, content, fData)
}

func TestOutputFile_LongFileName(t *testing.T) {
	directory := "./test_output1"
	defer func() {
		// remove UT output
		assert.NoError(t, os.RemoveAll(directory))
	}()

	fileName := "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest.html"
	content := []byte("test")

	o, err := NewOutputer(conf.OutputerConf{directory, ".*.(htm|html)$"})
	assert.NoError(t, err)

	err = o.OutputFile(fileName, content)
	assert.NoError(t, err)

	fp := path.Join(directory, hashLongFileName(fileName))
	fData, err := ioutil.ReadFile(fp)
	assert.NoError(t, err)
	assert.Equal(t, content, fData)
}

func TestOutputFile_FileNameNotMatch(t *testing.T) {
	directory := "./test_output1"
	defer func() {
		// remove UT output
		assert.NoError(t, os.RemoveAll(directory))
	}()

	fileName := "notMatchFileName"
	content := []byte("test")

	o, err := NewOutputer(conf.OutputerConf{directory, ".*.(htm|html)$"})
	assert.NoError(t, err)

	err = o.OutputFile(fileName, content)
	assert.NoError(t, err)

	fp := path.Join(directory, hashLongFileName(fileName))
	_, err = ioutil.ReadFile(fp)
	assert.True(t, strings.Contains(err.Error(), "no such file or directory"))
}
