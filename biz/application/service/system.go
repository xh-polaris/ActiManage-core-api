package service

import (
	"context"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	genbasic "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/basic"
	gensystem "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/system"
	"github.com/xh-polaris/ActiManage-core-api/biz/adaptor"
	"github.com/xh-polaris/ActiManage-core-api/biz/application/dto/core_api"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/config"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/consts"
	rpcsystem "github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/rpc/system"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/util"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/util/log"
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
	SystemRpc rpcsystem.IActiManageSystem
	Config    config.Config
}

var SystemServiceSet = wire.NewSet(
	wire.Struct(new(SystemService), "*"),
	wire.Bind(new(ISystemService), new(*SystemService)),
)

func (s SystemService) SystemLogin(ctx context.Context, req *core_api.SystemLoginReq) (resp *core_api.SystemLoginResp, err error) {
	loginResp, err := s.SystemRpc.AdminLogin(ctx, &gensystem.AdminLoginReq{
		AuthId:   req.AuthId,
		Password: req.Password,
	})
	if err != nil || loginResp.UserId == "" {
		log.Error("Login error: ", err)
		return nil, consts.ErrLogin
	}
	token, expire, err := util.GenerateJwtToken(s.Config.Auth.SecretKey, s.Config.Auth.AccessExpire, loginResp.UserId)
	if err != nil {
		log.Error("Generate jwt token error: ", err)
		return nil, consts.ErrLogin
	}
	resp = &core_api.SystemLoginResp{
		Code:        0,
		Msg:         "success",
		Id:          loginResp.UserId,
		AccessToken: token,
		ExpireTime:  expire,
	}
	return resp, nil
}

func (s SystemService) SystemListMerchant(ctx context.Context, req *core_api.SystemListMerchantsReq) (resp *core_api.SystemListMerchantsResp, err error) {
	response, err := s.SystemRpc.ListMerchants(ctx, &gensystem.ListMerchantsReq{
		Paging: &genbasic.Paging{
			Page:  req.Paging.Page,
			Limit: req.Paging.Limit,
		},
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	dtos := make([]*core_api.SystemListMerchantsResp_Item, len(response.Merchants))
	for _, merchant := range response.Merchants {
		dto := &core_api.SystemListMerchantsResp_Item{
			Id:   merchant.Id,
			Name: merchant.Name,
		}
		dtos = append(dtos, dto)
	}
	resp = &core_api.SystemListMerchantsResp{
		Code:      0,
		Msg:       "success",
		Merchants: dtos,
		Total:     response.Total,
	}
	return resp, nil
}

func (s SystemService) SystemGetMerchant(ctx context.Context, req *core_api.SystemGetMerchantReq) (resp *core_api.SystemGetMerchantResp, err error) {
	response, err := s.SystemRpc.GetMerchantInfo(ctx, &gensystem.GetMerchantInfoReq{
		Id: req.Id,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	resp = &core_api.SystemGetMerchantResp{}
	err = copier.Copy(&resp.Openings, &resp.Openings)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(&resp.Openings, &resp.Location)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Msg = "success"
	return resp, nil
}

func (s SystemService) SystemCreateMerchant(ctx context.Context, req *core_api.SystemCreateMerchantReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	rpcReq := &gensystem.CreateMerchantReq{
		AdminId:     userId,
		Name:        req.Name,
		Phone:       req.Phone,
		Description: req.Description,
		Licences:    req.Licences,
		Logo:        req.Logo,
	}
	err = copier.Copy(&rpcReq.Openings, &req.Openings)
	if err != nil {
		return nil, consts.ErrInvalidParameter
	}
	err = copier.Copy(&rpcReq.Location, &req.Location)
	if err != nil {
		return nil, consts.ErrInvalidParameter
	}
	response, err := s.SystemRpc.CreateMerchant(ctx, rpcReq)
	if err != nil || response.Code != 0 {
		return nil, consts.ErrCreate
	}
	return util.SuccessResp("商家创建成功")
}

func (s SystemService) SystemUpdateMerchant(ctx context.Context, req *core_api.SystemUpdateMerchantReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	rpcReq := &gensystem.UpdateMerchantReq{
		AdminId: userId,
		Merchant: &gensystem.Merchant{
			Id:          req.MerchantId,
			Name:        req.Name,
			Logo:        req.Logo,
			Phone:       req.Phone,
			Description: req.Description,
			Licences:    req.Licences,
			Status:      req.Status,
		},
	}
	err = copier.Copy(&rpcReq.Merchant.Openings, &req.Openings)
	if err != nil {
		return nil, consts.ErrInvalidParameter
	}
	err = copier.Copy(&rpcReq.Merchant.Location, &req.Location)
	if err != nil {
		return nil, consts.ErrInvalidParameter
	}
	response, err := s.SystemRpc.UpdateMerchant(ctx, rpcReq)
	if err != nil || response.Code != 0 {
		return nil, consts.ErrCreate
	}
	return util.SuccessResp("商家更新成功")
}

func (s SystemService) SystemGetDashboard(ctx context.Context, req *core_api.SystemGetDashboardReq) (resp *core_api.SystemGetDashboardResp, err error) {
	//TODO 根据图表计算数据
	panic("implement me")
}

func (s SystemService) SystemGetOverallDashboard(ctx context.Context, req *core_api.SystemGetOverallDashboardReq) (resp *core_api.SystemGetOverallDashboardResp, err error) {
	//TODO 根据图表计算数据
	panic("implement me")
}
