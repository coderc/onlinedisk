package file_handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeChunk(t *testing.T) {
	filePath := "./test_data/"
	err := mergeChunk(filePath, "_", 5)
	assert.Nil(t, err)
	testFile, err := os.Create("./test_data/test_file")
	assert.Nil(t, err)
	for i := 0; i < 5; i++ {
		file, err := os.Open(fmt.Sprintf("%s_%d", filePath, i))
		assert.Nil(t, err)
		defer file.Close()

		_, err = io.Copy(testFile, file)
		assert.Nil(t, err)
	}

	defer testFile.Close()
	ok, err := compareFiles("./test_data/test_file", "./test_data/_")
	assert.Nil(t, err)
	assert.True(t, ok)
}

func compareFiles(file1, file2 string) (bool, error) {
	content1, err := ioutil.ReadFile(file1)
	if err != nil {
		return false, err
	}

	content2, err := ioutil.ReadFile(file2)
	if err != nil {
		return false, err
	}

	return string(content1) == string(content2), nil
}
