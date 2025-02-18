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
	ErrInvalidParameter  = NewErrno(400, errors.New("invalid parameter"))
	ErrNotAuthentication = NewErrno(403, errors.New("not authentication"))
)

// 业务提示，code大于1000
var (
	ErrSignUp      = NewErrno(1000, errors.New("注册失败"))
	ErrLogin       = NewErrno(1001, errors.New("登录失败"))
	ErrSetPassword = NewErrno(1002, errors.New("密码修改失败"))
	ErrCreate      = NewErrno(2001, errors.New("创建失败"))
	ErrDelete      = NewErrno(2002, errors.New("删除失败"))
	ErrUpdate      = NewErrno(2003, errors.New("更新失败"))
	ErrCall        = NewErrno(3001, errors.New("调用模型失败"))
	ErrSend        = NewErrno(3002, errors.New("验证码发送失败"))
	ErrVerifyCode  = NewErrno(3003, errors.New("验证码错误"))
)
