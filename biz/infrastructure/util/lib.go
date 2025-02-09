package util

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/xh-polaris/ActiManage-core-api/biz/application/dto/core_api"
	"reflect"
	"time"
)

type Resp interface {
	GetCode() int64
	GetMsg() string
}

func SuccessResp(msg string) (*core_api.Response, error) {
	return &core_api.Response{
		Code: 0,
		Msg:  msg,
	}, nil
}

func FailResp(code int64, msg string) *core_api.Response {
	return &core_api.Response{
		Code: code,
		Msg:  msg,
	}
}

func FailRespStruct(resp *Resp, code int64, msg string) any {
	// 确保传入的是指针类型
	v := reflect.ValueOf(resp)

	// 获取结构体中的 Code 和 Msg 字段
	codeField := v.Elem().FieldByName("Code")
	msgField := v.Elem().FieldByName("Msg")

	// 修改 Code 和 Msg 字段
	if codeField.IsValid() && codeField.CanSet() {
		codeField.SetInt(code)
	}
	if msgField.IsValid() && msgField.CanSet() {
		msgField.SetString(msg)
	}

	return resp
}
func GenerateJwtToken(secret string, expire int64, userId string) (string, int64, error) {
	key, err := jwt.ParseECPrivateKeyFromPEM([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	iat := time.Now().Unix()
	exp := iat + expire
	claims := make(jwt.MapClaims)
	claims["exp"] = exp
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodES256)
	token.Claims = claims
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", 0, err
	}
	return tokenString, exp, nil
}
