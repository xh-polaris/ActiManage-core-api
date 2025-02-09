package adaptor

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	hertz "github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/util/log"
	bizerrors "github.com/xh-polaris/gopkg/errors"
	"github.com/xh-polaris/gopkg/util"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/status"
	"reflect"
)

var _ propagation.TextMapCarrier = &headerProvider{}

type headerProvider struct {
	headers *protocol.ResponseHeader
}

// Get a value from metadata by key
func (m *headerProvider) Get(key string) string {
	return m.headers.Get(key)
}

// Set a value to metadata by k/v
func (m *headerProvider) Set(key, value string) {
	m.headers.Set(key, value)
}

// Keys Iteratively get all keys of metadata
func (m *headerProvider) Keys() []string {
	out := make([]string, 0)

	m.headers.VisitAll(func(key, value []byte) {
		out = append(out, string(key))
	})

	return out
}

func PostProcess(ctx context.Context, c *app.RequestContext, req, resp any, err error) {
	log.CtxInfo(ctx, "[%s] req=%s, resp=%s, err=%v", c.Path(), util.JSONF(req), util.JSONF(resp), err)
	b3.New().Inject(ctx, &headerProvider{headers: &c.Response.Header})

	if err == nil { // 无错，返回响应
		response := makeResponse(resp)
		c.JSON(hertz.StatusOK, response)
	} else if s, ok := status.FromError(err); ok {
		StatusCode := hertz.StatusOK
		if s.Code() < 1000 {
			StatusCode = int(s.Code())
		}
		c.JSON(StatusCode, &bizerrors.BizError{
			Code: uint32(s.Code()),
			Msg:  s.Message(),
		})
	} else {
		log.CtxError(ctx, "internal error, err=%s", err.Error())
		code := hertz.StatusInternalServerError
		c.String(code, err.Error())
	}
}

func makeResponse(resp any) map[string]any {
	v := reflect.ValueOf(resp)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.Kind() == reflect.Struct {
			// 构建返回数据
			response := make(map[string]any, 3)
			response["code"] = v.FieldByName("Code").Int()
			response["msg"] = v.FieldByName("Msg").String()

			if v.NumField() == 2 {
				return response
			}

			data := make(map[string]interface{}, v.NumField()-2)
			for i := 0; i < v.NumField(); i++ {
				field := v.Type().Field(i)
				fieldValue := v.Field(i)
				jsonTag := field.Tag.Get("json")
				if field.Name == "code" || field.Name == "msg" {
					continue
				}
				// 获取 json 标签名，空值则用字段名
				if jsonTag == "" {
					jsonTag = field.Name
				}

				// 过滤零值字段，避免返回不必要的字段
				if !fieldValue.IsZero() {
					data[jsonTag] = fieldValue.Interface()
				}
			}
			response["data"] = data
			return response
		}
	}
	return nil
}
