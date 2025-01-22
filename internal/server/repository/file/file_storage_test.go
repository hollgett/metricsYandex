package file

import (
	"os"
	"testing"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
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

	fs, err := NewFileStorage(tempfile.Name())
	require.NoError(t, err, "create file storage error")

	err = fs.Update(data)
	require.NoError(t, err, "update data error")

	got, err := os.ReadFile(tempfile.Name())
	require.NoError(t, err, "read file error")
	assert.Equal(t, expected, string(got), "actual data not equal")

	err = fs.Close()
	require.NoError(t, err, "close file storage error")
}

func Test_fileStorage_Load(t *testing.T) {
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

	fs, err := NewFileStorage(fileTemp.Name())

	require.NoError(t, err, "create file storage error")

	loadData, err := fs.Load()
	require.NoError(t, err, "load file storage error")
	assert.Equal(t, expected, loadData, "not equal data")

	err = fs.Close()
	require.NoError(t, err, "close file storage error")

}
