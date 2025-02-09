package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/ActiManage-core-api/biz/application/dto/core_api"
)

type IUserService interface {
	Login(ctx context.Context, req *core_api.LoginReq) (resp *core_api.LoginResp, err error)
	SignUp(ctx context.Context, req *core_api.SignUpReq) (resp *core_api.SignUpResp, err error)
	GetSetting(ctx context.Context, req *core_api.GetSettingReq) (resp *core_api.GetActivityResp, err error)
	ListActivities(ctx context.Context, req *core_api.ListActivitiesReq) (resp *core_api.ListActivitiesResp, err error)
	GetActivity(ctx context.Context, req *core_api.GetActivityReq) (resp *core_api.GetActivityResp, err error)
	DoFavorite(ctx context.Context, req *core_api.DoFavoriteReq) (resp *core_api.Response, err error)
	CancelFavorite(ctx context.Context, req *core_api.CancelFavoriteReq) (resp *core_api.Response, err error)
	CreateBooking(ctx context.Context, req *core_api.CreateBookingReq) (resp *core_api.Response, err error)
	CancelBookRecord(ctx context.Context, req *core_api.CancelBookRecordReq) (resp *core_api.Response, err error)
	ListActivitiesByBookRecords(ctx context.Context, req *core_api.ListBookRecordsReq) (resp *core_api.ListBookRecordsResp, err error)
	ListReservers(ctx context.Context, req *core_api.ListReserversReq) (resp *core_api.ListReserversResp, err error)
	CreateReserver(ctx context.Context, req *core_api.CreateReserverReq) (resp *core_api.Reserver, err error)
	DeleteReserver(ctx context.Context, req *core_api.DeleteReserverReq) (resp *core_api.Reserver, err error)
	GetUserInfo(ctx context.Context, req *core_api.GetUserInfoReq) (resp *core_api.GetUserInfoResp, err error)
	UpdateUserInfo(ctx context.Context, req *core_api.UpdateUserInfoReq) (resp *core_api.Response, err error)
	UpdateNotice(ctx context.Context, req *core_api.UpdateNoticeReq) (resp *core_api.Response, err error)
}

type UserService struct {
}

var UserServiceSet = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),
)

func (u UserService) Login(ctx context.Context, req *core_api.LoginReq) (resp *core_api.LoginResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) SignUp(ctx context.Context, req *core_api.SignUpReq) (resp *core_api.SignUpResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) GetSetting(ctx context.Context, req *core_api.GetSettingReq) (resp *core_api.GetActivityResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) ListActivities(ctx context.Context, req *core_api.ListActivitiesReq) (resp *core_api.ListActivitiesResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) GetActivity(ctx context.Context, req *core_api.GetActivityReq) (resp *core_api.GetActivityResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) DoFavorite(ctx context.Context, req *core_api.DoFavoriteReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) CancelFavorite(ctx context.Context, req *core_api.CancelFavoriteReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) CreateBooking(ctx context.Context, req *core_api.CreateBookingReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) CancelBookRecord(ctx context.Context, req *core_api.CancelBookRecordReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) ListActivitiesByBookRecords(ctx context.Context, req *core_api.ListBookRecordsReq) (resp *core_api.ListBookRecordsResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) ListReservers(ctx context.Context, req *core_api.ListReserversReq) (resp *core_api.ListReserversResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) CreateReserver(ctx context.Context, req *core_api.CreateReserverReq) (resp *core_api.Reserver, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) DeleteReserver(ctx context.Context, req *core_api.DeleteReserverReq) (resp *core_api.Reserver, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) GetUserInfo(ctx context.Context, req *core_api.GetUserInfoReq) (resp *core_api.GetUserInfoResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) UpdateUserInfo(ctx context.Context, req *core_api.UpdateUserInfoReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) UpdateNotice(ctx context.Context, req *core_api.UpdateNoticeReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}
