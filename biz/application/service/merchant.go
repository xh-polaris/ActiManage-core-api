package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/ActiManage-core-api/biz/application/dto/core_api"
)

type IMerchantService interface {
	MerchantListActivities(ctx context.Context, c *core_api.MerchantListActivitiesReq) (resp *core_api.MerchantListActivitiesResp, err error)
	MerchantCreateActivity(ctx context.Context, c *core_api.MerchantCreateActivityReq) (resp *core_api.Response, err error)
	MerchantDeleteActivity(ctx context.Context, c *core_api.MerchantDeleteActivityReq) (resp *core_api.Response, err error)
	MerchantTopActivity(ctx context.Context, c *core_api.MerchantTopActivityReq) (resp *core_api.Response, err error)
	MerchantLogin(ctx context.Context, c *core_api.MerchantLoginReq) (resp *core_api.MerchantLoginResp, err error)
	MerchantGetSetting(ctx context.Context, c *core_api.MerchantGetSettingReq) (resp *core_api.MerchantGetSettingResp, err error)
	MerchantUpdateSetting(ctx context.Context, c *core_api.MerchantUpdateSettingReq) (resp *core_api.Response, err error)
	MerchantGetBookRecords(ctx context.Context, c *core_api.MerchantListBookRecordsReq) (resp *core_api.MerchantListBookRecordsResp, err error)
	MerchantUpdateInfo(ctx context.Context, c *core_api.MerchantUpdateInfoReq) (resp *core_api.Response, err error)
	MerchantGetInfo(ctx context.Context, c *core_api.MerchantGetInfoReq) (resp *core_api.MerchantGetInfoResp, err error)
	MerchantSetPassword(ctx context.Context, c *core_api.MerchantSetPasswordReq) (resp *core_api.Response, err error)
}

type MerchantService struct {
}

var MerchantServiceSet = wire.NewSet(
	wire.Struct(new(MerchantService), "*"),
	wire.Bind(new(IMerchantService), new(*MerchantService)),
)

func (s *MerchantService) MerchantListActivities(ctx context.Context, req *core_api.MerchantListActivitiesReq) (resp *core_api.MerchantListActivitiesResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchantService) MerchantCreateActivity(ctx context.Context, c *core_api.MerchantCreateActivityReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchantService) MerchantDeleteActivity(ctx context.Context, c *core_api.MerchantDeleteActivityReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchantService) MerchantTopActivity(ctx context.Context, c *core_api.MerchantTopActivityReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchantService) MerchantLogin(ctx context.Context, c *core_api.MerchantLoginReq) (resp *core_api.MerchantLoginResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchantService) MerchantGetSetting(ctx context.Context, c *core_api.MerchantGetSettingReq) (resp *core_api.MerchantGetSettingResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchantService) MerchantUpdateSetting(ctx context.Context, c *core_api.MerchantUpdateSettingReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchantService) MerchantGetBookRecords(ctx context.Context, c *core_api.MerchantListBookRecordsReq) (resp *core_api.MerchantListBookRecordsResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchantService) MerchantUpdateInfo(ctx context.Context, c *core_api.MerchantUpdateInfoReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchantService) MerchantGetInfo(ctx context.Context, c *core_api.MerchantGetInfoReq) (resp *core_api.MerchantGetInfoResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchantService) MerchantSetPassword(ctx context.Context, c *core_api.MerchantSetPasswordReq) (resp *core_api.Response, err error) {
	//TODO implement me
	panic("implement me")
}
