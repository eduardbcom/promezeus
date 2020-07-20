package promquery

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// Config structure
type Config struct {
	PromAPI string
}

// Service structure
type Service struct {
	config *Config
	client api.Client
}

// New creates a new instance of promquery.Service
func New(config *Config) (*Service, error) {
	client, err := api.NewClient(api.Config{
		Address: config.PromAPI,
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		config: config,
		client: client,
	}, nil
}

// Query fetches data from prometheus using PromQL
func (r *Service) Query(ctx context.Context, query string) (map[string]interface{}, error) {
	v1api := v1.NewAPI(r.client)

	resp, warnings, err := v1api.Query(ctx, query, time.Time{})

	if err != nil {
		return nil, err
	}

	if len(warnings) > 0 {
		return nil, errors.New(warnings[0])
	}

	var res []map[string]interface{}
	result, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(result, &res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return convertResultToResponse(res), nil
}

func convertResultToResponse(res []map[string]interface{}) map[string]interface{} {
	response := make(map[string]interface{})

	for _, resItem := range res {
		valI := resItem["value"]
		metricI := resItem["metric"]

		valS := valI.([]interface{})

		resObj := map[string]interface{}{
			"timestamp": valS[0].(float64),
			"value":     valS[1].(string),
		}

		metric := metricI.(map[string]interface{})
		val, ok := metric["id"]
		if ok {
			response[val.(string)] = resObj

			continue
		}

		response = resObj
	}

	return response
}
