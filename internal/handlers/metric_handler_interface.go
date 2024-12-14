package handlers

//go:generate mockgen -source=metric_handler_interface.go -destination=../mock/metric_handler.go -package=mock
type MetricHandler interface {
	CollectingMetric(requestParam []string) error
	GetMetric(requestParam []string) (string, error)
	GetMetricAll() (string, error)
}


