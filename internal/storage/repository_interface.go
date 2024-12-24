package storage

//go:generate mockgen -source=repository_interface.go -destination=../mock/repository.go -package=mock
type Repositories interface {
	UpdateGauge(nameMetric string, val float64) error
	AddCounter(nameMetric string, val int64) error
	GetMetricGauge(nameM string) (float64, error)
	GetMetricCounter(nameM string) (int64, error)
	GetMetricAll() (map[string]string, error)
}
