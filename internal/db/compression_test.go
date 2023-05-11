package db

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := logging.SetDefaultLogger("CONSOLE")
	if err != nil {
		panic(err)
	}

	logger = logging.GetLogger()

	code := m.Run()
	os.Exit(code)
}

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
