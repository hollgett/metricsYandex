package repository

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hollgett/metricsYandex.git/internal/server/logger"
	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/repository/file"
	"github.com/hollgett/metricsYandex.git/internal/server/repository/memory"
)

var (
	ErrMetric = errors.New("unknown metrics")
)

type compositeStorage struct {
	fileStorage   file.File
	memoryStorage memory.Memory
}

func NewCompositeStorage(path string, restore bool) (Repository, error) {
	compStor := &compositeStorage{
		memoryStorage: memory.NewMemoryStorage(),
	}
	validPath := path != ""

	if validPath {
		file, err := file.NewFileStorage(path)
		if err != nil {
			return nil, fmt.Errorf("file storage create error: %w", err)
		}
		compStor.fileStorage = file
	}

	if restore && validPath {
		fStorage, err := compStor.fileStorage.Load()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("file storage load error: %w", err)
		}
		for _, v := range fStorage {
			if err := compStor.Save(v); err != nil {
				return nil, fmt.Errorf("file storage save error: %w", err)
			}
		}
		logger.LogAny("restore data", "len", len(fStorage))
	}
	return compStor, nil
}

func (cs *compositeStorage) Save(data models.Metrics) error {
	switch data.MType {
	case "gauge":
		if err := cs.memoryStorage.SetGauge(data.ID, *data.Value); err != nil {
			return fmt.Errorf("set gauge err: %w", err)
		}
	case "counter":
		if err := cs.memoryStorage.AddCounter(data.ID, *data.Delta); err != nil {
			return fmt.Errorf("add counter err: %w", err)
		}
	default:
		return ErrMetric
	}
	return nil
}

func (cs *compositeStorage) Get(metric *models.Metrics) error {
	switch metric.MType {
	case "gauge":
		val, err := cs.memoryStorage.GetGauge(metric.ID)
		if err != nil {
			return fmt.Errorf("get gauge err: %w", err)
		}
		metric.Value = &val
	case "counter":
		val, err := cs.memoryStorage.GetCounter(metric.ID)
		if err != nil {
			return fmt.Errorf("get counter err: %w", err)
		}
		metric.Delta = &val
	default:
		return ErrMetric
	}
	return nil
}

func (cs *compositeStorage) GetAll() ([]models.Metrics, error) {
	listData := cs.memoryStorage.GetAll()
	if len(listData) == 0 {
		return nil, errors.New("data base doesn't have metrics")
	}
	return listData, nil
}

func (cs *compositeStorage) Close() error {
	dStorage := cs.memoryStorage.GetAll()
	if err := cs.fileStorage.Update(dStorage); err != nil {
		logger.LogErr("close update data", err)
	}
	return cs.fileStorage.Close()
}

func (cs *compositeStorage) UpdateTicker(interval int, done chan bool) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-done:
			return
		default:
			dStorage := cs.memoryStorage.GetAll()
			if err := cs.fileStorage.Update(dStorage); err != nil {
				logger.LogErr("update file", err)
			}
			logger.LogAny("update file", "complete", len(dStorage))
		}
	}
}

func (cs *compositeStorage) UpdateSyncMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "update") {
			next.ServeHTTP(w, r)
			dStorage := cs.memoryStorage.GetAll()
			if err := cs.fileStorage.Update(dStorage); err != nil {
				logger.LogErr("sync update", err)
			}
			logger.LogAny("update file", "complete", len(dStorage))
		}

	})
}
