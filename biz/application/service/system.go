package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/ActiManage-core-api/biz/application/dto/core_api"
)

type ISystemService interface {
	SystemLogin(ctx context.Context, req *core_api.SystemLoginReq) (resp *core_api.SystemLoginResp, err error)
	SystemListMerchant(ctx context.Context, req *core_api.SystemListMerchantsReq) (resp *core_api.SystemListMerchantsResp, err error)
	SystemGetMerchant(ctx context.Context, req *core_api.SystemGetMerchantReq) (resp *core_api.SystemGetMerchantResp, err error)
	SystemCreateMerchant(ctx context.Context, req *core_api.SystemCreateMerchantReq) (resp *core_api.Response, err error)
	SystemUpdateMerchant(ctx context.Context, req *core_api.SystemUpdateMerchantReq) (resp *core_api.Response, err error)
	SystemGetDashboard(ctx context.Context, req *core_api.SystemGetDashboardReq) (resp *core_api.SystemGetDashboardResp, err error)
	SystemGetOverallDashboard(ctx context.Context, req *core_api.SystemGetOverallDashboardReq) (resp *core_api.SystemGetOverallDashboardResp, err error)
}

type SystemService struct {
}

var SystemServiceSet = wire.NewSet(
	wire.Struct(new(SystemService), "*"),
	wire.Bind(new(ISystemService), new(*SystemService)),
)

func (s SystemService) SystemLogin(ctx context.Context, req *core_api.SystemLoginReq) (resp *core_api.SystemLoginResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemService) SystemListMerchant(ctx context.Context, req *core_api.SystemListMerchantsReq) (resp *core_api.SystemListMerchantsResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemService) SystemGetMerchant(ctx context.Context, req *core_api.SystemGetMerchantReq) (resp *core_api.SystemGetMerchantResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemService) SystemCreateMerchant(ctx context.Context, req *core_api.SystemCreateMerchantReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemService) SystemUpdateMerchant(ctx context.Context, req *core_api.SystemUpdateMerchantReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemService) SystemGetDashboard(ctx context.Context, req *core_api.SystemGetDashboardReq) (resp *core_api.SystemGetDashboardResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemService) SystemGetOverallDashboard(ctx context.Context, req *core_api.SystemGetOverallDashboardReq) (resp *core_api.SystemGetOverallDashboardResp, err error) {
	//TODO implement me
	panic("implement me")
}
