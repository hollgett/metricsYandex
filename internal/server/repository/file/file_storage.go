package file

import (
	"io"
	"os"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/utils"
)

type File interface {
	Update(dataStor []models.Metrics) error
	Load() ([]models.Metrics, error)
	Close() error
}


type FileStorage struct {
	file *os.File
}

func NewFileStorage(dir string) (File, error) {
	file, err := os.OpenFile(dir, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &FileStorage{
		file: file,
	}, nil
}

func (fs *FileStorage) Close() error {
	return fs.file.Close()
}

func (fs *FileStorage) Update(dataStor []models.Metrics) error {
	if err := fs.file.Truncate(0); err != nil {
		return err
	}
	if _, err := fs.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := utils.EncoderJSON(fs.file, dataStor); err != nil {
		return err
	}
	if err := fs.file.Sync(); err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) Load() ([]models.Metrics, error) {
	var data []models.Metrics
	if err := utils.DecoderJSON(fs.file, &data); err != nil {
		return nil, err
	}
	return data, nil
}
