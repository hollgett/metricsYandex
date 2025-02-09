package api

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hollgett/metricsYandex.git/internal/agent/config"
	"github.com/hollgett/metricsYandex.git/internal/agent/logger"
	"github.com/hollgett/metricsYandex.git/internal/agent/models"
	"github.com/hollgett/metricsYandex.git/internal/agent/utils"
	"go.uber.org/zap"
)

type Client struct {
	Client *resty.Client
}

func NewClientResty(header, value string, debug bool) *Client {
	client := resty.New().
		SetBaseURL(`http://`+config.AgentConfig.Addr).
		SetHeader(header, value).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetDebug(debug)
	return &Client{
		Client: client,
	}
}

func (c *Client) request(metric models.Metrics) (*resty.Response, error) {
	data, err := utils.Marshal(metric)
	if err != nil {
		return nil, fmt.Errorf("encode json: %w", err)
	}
	data, err = utils.CompressData(data)
	if err != nil {
		return nil, fmt.Errorf("compress: %w", err)
	}

	return c.Client.R().
		SetBody(data).
		Post(`/update/`)
}

func (c *Client) SendMetricsJSON(dataMetric []models.Metrics) error {
	for _, metric := range dataMetric {
		resp, err := c.request(metric)
		if err != nil {
			return fmt.Errorf("request: %w", err)
		}
		logger.Log.Info("request", zap.Any("value", metric), zap.String("status", resp.Status()))
	}
	return nil
}
