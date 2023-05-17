package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompressDecompress(t *testing.T) {
	input := "Some string input"

	compressed, err := compress(input)
	if err != nil {
		t.Error(err)
	}

	output, err := decompress(compressed)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, input, output)
}
