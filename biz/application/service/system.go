package service

import (
	"context"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	genbasic "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/basic"
	gensystem "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/system"
	genuser "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/user"
	"github.com/xh-polaris/ActiManage-core-api/biz/adaptor"
	"github.com/xh-polaris/ActiManage-core-api/biz/application/dto/core_api"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/config"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/consts"
	rpcsystem "github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/rpc/system"
	rpcuser "github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/rpc/user"
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
	UserRpc   rpcuser.IActiManageUser
	Config    *config.Config
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
	dtos := make([]*core_api.SystemListMerchantsResp_Item, 0, len(response.Merchants))
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
	resp = &core_api.SystemGetMerchantResp{
		Code:      0,
		Msg:       "success",
		Location:  &core_api.Location{},
		Openings:  make([]*core_api.Opening, 0, len(response.Openings)),
		Establish: response.Establish,
		Capital:   response.Capital,
	}
	err = copier.Copy(&resp, &response)
	if err != nil {
		return nil, err
	}
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
		Location:    &gensystem.Location{},
		Openings:    make([]*gensystem.Opening, 0, len(req.Openings)),
		Logo:        req.Logo,
		Establish:   req.Establish,
		Capital:     req.Capital,
	}
	err = copier.Copy(&rpcReq, &req)
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
			Id:        req.MerchantId,
			Establish: req.Establish,
			Capital:   req.Capital,
			Location:  &gensystem.Location{},
			Openings:  make([]*gensystem.Opening, 0, len(req.Openings)),
		},
	}
	err = copier.Copy(&rpcReq.Merchant, &req)
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
	// 获取访问量折线图
	viewLineResp, err := s.UserRpc.GetViewDataByMerchant(ctx, &genuser.GetViewDataByMerchantReq{
		Number:     req.ViewDataNumber,
		MerchantId: req.Id,
	})
	if err != nil {
		return nil, err
	}
	views := make([]*core_api.SystemGetDashboardResp_ViewItem, 0, len(viewLineResp.Items))
	for _, v := range viewLineResp.Items {
		views = append(views, &core_api.SystemGetDashboardResp_ViewItem{
			Number: v.Number,
			Time:   v.Time,
		})
	}
	// 根据预约量获取活动id
	activityIdResp, err := s.UserRpc.ListActivityIdByBookRecordRank(ctx, &genuser.ListActivityIdsByBookRecordRankReq{Number: req.ActivityByBookNumber, MerchantId: req.Id})
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(activityIdResp.Items))
	for _, item := range activityIdResp.Items {
		ids = append(ids, item.Id)
	}
	activityResp, err := s.SystemRpc.ListActivityByActivityId(ctx, &gensystem.ListActivitiesByActivityIdReq{Ids: ids})
	if err != nil {
		return nil, err
	}
	activitiesByBookRecordNumber := make([]*core_api.SystemGetDashboardResp_ActivityItem, 0, len(activityResp.Items))
	for i, item := range activityResp.Items {
		activitiesByBookRecordNumber = append(activitiesByBookRecordNumber, &core_api.SystemGetDashboardResp_ActivityItem{
			Id:     item.Id,
			Name:   item.Name,
			Number: activityIdResp.Items[i].Number,
		})
	}

	resp = &core_api.SystemGetDashboardResp{
		Code:                 0,
		Msg:                  "success",
		ViewData:             views,
		ActivityByBookNumber: activitiesByBookRecordNumber,
	}
	return resp, nil
}

func (s SystemService) SystemGetOverallDashboard(ctx context.Context, req *core_api.SystemGetOverallDashboardReq) (resp *core_api.SystemGetOverallDashboardResp, err error) {
	// 商家总数折线图
	merchantTotal, err := s.SystemRpc.GetMerchantTotalData(ctx, &gensystem.GetMerchantTotalDataReq{
		Number: req.LineNumber,
	})
	if err != nil {
		return nil, err
	}
	totals := make([]*core_api.SystemGetOverallDashboardResp_LineItem, 0, len(merchantTotal.Items))
	for _, v := range merchantTotal.Items {
		totals = append(totals, &core_api.SystemGetOverallDashboardResp_LineItem{
			Number: v.Number,
			Time:   v.Time,
		})
	}
	// 商家访问量排名
	merchantIdByViewResp, err := s.UserRpc.ListMerchantIdByViewRank(ctx, &genuser.ListMerchantIdsByViewRankReq{Number: req.MerchantByViewRankNumber})
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(merchantIdByViewResp.Items))
	for _, item := range merchantIdByViewResp.Items {
		ids = append(ids, item.Id)
	}
	merchantByViewResp, err := s.SystemRpc.ListMerchantByMerchantId(ctx, &gensystem.ListMerchantsByMerchantIdReq{Ids: ids})
	if err != nil {
		return nil, err
	}
	merchantByView := make([]*core_api.SystemGetOverallDashboardResp_MerchantItem, 0, len(merchantByViewResp.Items))
	for i, item := range merchantByViewResp.Items {
		merchantByView = append(merchantByView, &core_api.SystemGetOverallDashboardResp_MerchantItem{
			Id:     item.Id,
			Name:   item.Name,
			Logo:   item.Logo,
			Number: merchantIdByViewResp.Items[i].Number,
		})
	}

	// 商家预约量排名
	merchantIdByBookRecordResp, err := s.UserRpc.ListMerchantIdByBookRecordRank(ctx, &genuser.ListMerchantIdsByBookRecordRankReq{Number: req.MerchantByBookRecordRankNumber})
	if err != nil {
		return nil, err
	}
	ids = make([]string, 0, len(merchantIdByBookRecordResp.Items))
	for _, item := range merchantIdByBookRecordResp.Items {
		ids = append(ids, item.Id)
	}
	merchantByBookRecordResp, err := s.SystemRpc.ListMerchantByMerchantId(ctx, &gensystem.ListMerchantsByMerchantIdReq{Ids: ids})
	if err != nil {
		return nil, err
	}
	merchantByBookRecord := make([]*core_api.SystemGetOverallDashboardResp_MerchantItem, 0, len(merchantByBookRecordResp.Items))
	for i, item := range merchantByBookRecordResp.Items {
		merchantByBookRecord = append(merchantByBookRecord, &core_api.SystemGetOverallDashboardResp_MerchantItem{
			Id:     item.Id,
			Name:   item.Name,
			Logo:   item.Logo,
			Number: merchantIdByBookRecordResp.Items[i].Number,
		})
	}

	// 商家活动量排名
	merchantByActivityResp, err := s.SystemRpc.ListMerchantByActivityNumber(ctx, &gensystem.ListMerchantsByActivityNumberReq{Number: req.MerchantByActivityNumberRankNumber})
	if err != nil {
		return nil, err
	}
	merchantByActivity := make([]*core_api.SystemGetOverallDashboardResp_MerchantItem, 0, len(merchantByActivityResp.Items))
	for _, item := range merchantByActivityResp.Items {
		merchantByActivity = append(merchantByActivity, &core_api.SystemGetOverallDashboardResp_MerchantItem{
			Id:     item.Id,
			Name:   item.Name,
			Logo:   item.Logo,
			Number: item.Number,
		})
	}

	resp = &core_api.SystemGetOverallDashboardResp{
		Code:                         0,
		Msg:                          "success",
		LineData:                     totals,
		MerchantByViewRank:           merchantByView,
		MerchantByBookRecordRank:     merchantByBookRecord,
		MerchantByActivityNumberRank: merchantByActivity,
	}
	return resp, nil
}

func (s SystemService) ResetMerchantPassword(ctx context.Context, req *core_api.ResetMerchantPasswordReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}

	response, err := s.SystemRpc.ResetMerchantPassword(ctx, &gensystem.ResetMerchantPasswordReq{
		MerchantId:  req.MerchantId,
		NewPassword: req.NewPassword,
		AdminId:     userId,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	resp = &core_api.Response{
		Code: 0,
		Msg:  "success",
	}
	return resp, nil
}
