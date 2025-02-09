package api

import (
	"context"
	"errors"
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

func (c *Client) SendMetricsJSON(ctx context.Context, dataMetric []models.Metrics, retryCount int, t time.Duration) error {
	ctxSend, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, metric := range dataMetric {
		if err := c.SendWithRetry(ctxSend, retryCount, t, metric); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) SendWithRetry(ctx context.Context, retryCount int, t time.Duration, metric models.Metrics) error {
	for i := 0; i <= retryCount; i++ {
		resp, err := c.request(metric)
		if err == nil && resp.IsSuccess() {
			logger.Log.Info("request", zap.Any("value", metric), zap.String("status", resp.Status()))
			break
		}
		select {
		case <-time.After(t):
			logger.Log.Info("request retry", zap.Error(err), zap.Int("count", i+1))
		case <-ctx.Done():
			logger.Log.Info("context cancel", zap.Error(ctx.Err()))
			return err
		}
		if i == 2 {
			return errors.New("limit retry exceeded")
		}
	}
	return nil
}
