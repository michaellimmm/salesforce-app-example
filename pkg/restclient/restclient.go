package restclient

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type (
	RestClient interface {
		OAuth
		UserInfo
	}

	restClient struct {
		logger *zap.Logger
		client *resty.Client
	}
)

func NewRestClient(logger *zap.Logger, client *resty.Client) RestClient {
	return &restClient{
		client: client,
		logger: logger,
	}
}
