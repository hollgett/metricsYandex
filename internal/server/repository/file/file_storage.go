package file

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/hollgett/metricsYandex.git/internal/server/logger"
	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/repository"
	"github.com/hollgett/metricsYandex.git/internal/server/repository/memory"
)

var (
	ErrMetric = errors.New("unknown metrics")
)

type FileStorage struct {
	file      *os.File
	updateInt int
	repository.Repository
	logger.Logger
}

func New(ctx context.Context, log logger.Logger, dir string, updateInt int, restore bool) (repository.Repository, error) {
	file, err := os.OpenFile(dir, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	fs := FileStorage{
		file:       file,
		updateInt:  updateInt,
		Repository: memory.New(),
		Logger:     log,
	}

	if restore {
		fs.restore()
	}

	if updateInt > 0 {
		go fs.updateTicker(ctx, updateInt)
	}

	return &fs, nil
}

func (fs *FileStorage) Close() error {
	return fs.file.Close()
}

func (fs *FileStorage) Save(data models.Metrics) error {
	fs.Repository.Save(data)
	if fs.updateInt == 0 {
		data, err := fs.GetAll()
		if err != nil {
			return err
		}
		if err := fs.update(data); err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileStorage) Ping(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return ctx.Err()
	}
	_, err := fs.file.Stat()
	if err != nil {
		fs.LogErr("ping file", err)
		return err
	}
	return nil
}

func (fs *FileStorage) Batch(ctx context.Context, metrics []models.Metrics) error {
	if err := ctx.Err(); err != nil {
		return ctx.Err()
	}
	for _, v := range metrics {
		if err := fs.Save(v); err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileStorage) updateTicker(ctx context.Context, updateInt int) {
	if err := ctx.Err(); err != nil {
		fs.Logger.LogErr("context", err)
		return
	}
	ticker := time.NewTicker(time.Duration(updateInt) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
			dStorage, err := fs.GetAll()
			if err != nil {
				fs.Logger.LogErr("storageget all", err)
			}
			if err := fs.update(dStorage); err != nil {
				fs.LogErr("update file", err)
			}
			fs.LogAny("update file", "complete", len(dStorage))
		}
	}
}

func (fs *FileStorage) update(dataStor []models.Metrics) error {
	if err := fs.file.Truncate(0); err != nil {
		return err
	}
	if _, err := fs.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := json.NewEncoder(fs.file).Encode(dataStor); err != nil {
		return err
	}
	if err := fs.file.Sync(); err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) restore() error {
	var data []models.Metrics
	if err := json.NewDecoder(fs.file).Decode(&data); err != nil {
		return err
	}
	for _, v := range data {
		fs.Repository.Save(v)
	}
	fs.LogAny("restore", "len", len(data))
	return nil
}
