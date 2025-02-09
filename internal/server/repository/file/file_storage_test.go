package file

import (
	"os"
	"testing"

	"github.com/hollgett/metricsYandex.git/internal/server/logger"
	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/repository/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_fileStorage_Update(t *testing.T) {
	var gauge = 54.3
	var counter int64 = 43
	data := []models.Metrics{
		{
			ID:    "gauge1",
			MType: "gauge",
			Value: &gauge,
		},
		{
			ID:    "counter1",
			MType: "counter",
			Delta: &counter,
		},
	}
	expected := `[{"id":"gauge1","type":"gauge","value":54.3},{"id":"counter1","type":"counter","delta":43}]
`

	tempfile, err := os.CreateTemp(t.TempDir(), "*.json")
	require.NoError(t, err, "create temp file error")
	tempfile.Close()
	defer os.RemoveAll(tempfile.Name())

	file, err := os.OpenFile(tempfile.Name(), os.O_RDWR|os.O_CREATE, 0666)
	require.NoError(t, err, "open temp file")
	fs := &FileStorage{
		file:       file,
		Repository: memory.New(),
	}
	err = fs.update(data)
	require.NoError(t, err, "update data error")

	got, err := os.ReadFile(tempfile.Name())
	require.NoError(t, err, "read file error")
	assert.Equal(t, expected, string(got), "actual data not equal")

	err = fs.Close()
	require.NoError(t, err, "close file storage error")
}

func Test_fileStorage_Restore(t *testing.T) {
	var gauge = 54.3
	var counter int64 = 43
	expected := []models.Metrics{
		{
			ID:    "gauge1",
			MType: "gauge",
			Value: &gauge,
		},
		{
			ID:    "counter1",
			MType: "counter",
			Delta: &counter,
		},
	}
	data := `[{"id":"gauge1","type":"gauge","value":54.3},{"id":"counter1","type":"counter","delta":43}]
`
	fileTemp, err := os.CreateTemp(t.TempDir(), "*.json")
	require.NoError(t, err, "create temp file error")
	_, err = fileTemp.Write([]byte(data))
	require.NoError(t, err, "write temp file error")
	err = fileTemp.Close()
	require.NoError(t, err, "close temp file error")
	defer os.Remove(fileTemp.Name())

	file, err := os.OpenFile(fileTemp.Name(), os.O_RDWR|os.O_CREATE, 0666)
	require.NoError(t, err, "open temp file")
	log, err := logger.New()
	require.NoError(t, err, "logger create")
	fs := &FileStorage{
		file:       file,
		Repository: memory.New(),
		Logger:     log,
	}
	require.NoError(t, err, "create file storage error")

	defer func() {
		err = fs.Close()
		require.NoError(t, err, "close file storage error")
	}()

	err = fs.restore()
	require.NoError(t, err, "load file storage error")
	dataGet, err := fs.GetAll()
	require.NoError(t, err, "get data storage error")
	assert.Equal(t, expected, dataGet, "not equal data")

}
