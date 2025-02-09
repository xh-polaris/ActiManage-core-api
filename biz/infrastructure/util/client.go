package util

import (
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/config"
	"net/http"
)

// HttpClient 是一个简单的 HTTP 客户端
type HttpClient struct {
	Client *http.Client
	Config *config.Config
}

// NewHttpClient 创建一个新的 HttpClient 实例
func NewHttpClient() *HttpClient {
	return &HttpClient{
		Client: &http.Client{},
	}
}
