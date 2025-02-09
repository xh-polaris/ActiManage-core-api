package consts

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Errno struct {
	err  error
	code int32
}

// GRPCStatus 实现 GRPCStatus 方法
func (en *Errno) GRPCStatus() *status.Status {
	return status.New(codes.Code(en.code), en.err.Error())
}

// 实现 Error 方法
func (en *Errno) Error() string {
	return en.err.Error()
}

func (en *Errno) Code() int32 {
	return en.code
}

// NewErrno 创建自定义错误
func NewErrno(code int32, err error) *Errno {
	return &Errno{
		err:  err,
		code: code,
	}
}

// 改变状态码错误，code小于1000
var (
	ErrNotAuthentication = NewErrno(403, errors.New("not authentication"))
)

// 业务提示，code大于1000
var ()
