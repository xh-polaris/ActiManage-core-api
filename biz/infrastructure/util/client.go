package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/config"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/consts"
)

// HttpClient 是一个简单的 HTTP 客户端
type HttpClient struct {
	Client *http.Client
	Config *config.Config
}

// NewHttpClient 创建一个新的 HttpClient 实例
func NewHttpClient() *HttpClient {
	// 创建一个禁用TLS验证的Transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &HttpClient{
		Client: &http.Client{Transport: transport},
	}
}

// CallGLM 调用GLM生成文案, text为原文案
func (c *HttpClient) CallGLM(text string, lang string) (map[string]interface{}, error) {
	header := make(map[string]string)
	header["Content-Type"] = consts.ContentTypeJson
	header["Charset"] = consts.CharSetUTF8
	header["Authorization"] = "Bearer " + config.GetConfig().GLMKey

	// 定义消息结构
	message := []map[string]interface{}{
		{
			"role": "user",
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": "这是一个活动的营销文案:" + text + "请帮我润色，要求有吸引力，能吸引顾客参加, 给我润色后的文案，不要有额外的输出。有适当的emoji和生动的语言，篇幅不要太短" +
						"并用日语输出",
				},
			},
		},
	}

	body := make(map[string]interface{})
	body["model"] = config.GetConfig().GLMModel
	body["messages"] = message
	body["thinking"] = map[string]interface{}{
		"type": "disabled",
	}

	resp, err := c.SendRequest(consts.Post, config.GetConfig().GLMUrl, header, body)
	fmt.Println("模型响应:", resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SendRequest 发送 HTTP 请求
func (c *HttpClient) SendRequest(method, url string, headers map[string]string, body interface{}) (map[string]interface{}, error) {
	// 将 body 序列化为 JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("请求体序列化失败: %w", err)
	}

	// 创建新的请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("关闭请求失败: %v", closeErr)
		}
	}()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("unexpected status code: %d, response body: %s", resp.StatusCode, responseBody)
		return nil, fmt.Errorf(errMsg)
	}

	// 反序列化响应体
	var responseMap map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseMap); err != nil {
		return nil, fmt.Errorf("反序列化响应失败: %w", err)
	}

	return responseMap, nil
}
