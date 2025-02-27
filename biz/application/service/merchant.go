package service

import (
	"context"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	genbasic "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/basic"
	gensystem "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/system"
	genuser "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/user"
	"github.com/xh-polaris/ActiManage-core-api/biz/adaptor"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/config"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/consts"
	rpcuser "github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/rpc/user"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/util"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/util/log"

	"github.com/xh-polaris/ActiManage-core-api/biz/application/dto/core_api"
	rpcsystem "github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/rpc/system"
)

type IMerchantService interface {
	MerchantListActivities(ctx context.Context, req *core_api.MerchantListActivitiesReq) (resp *core_api.MerchantListActivitiesResp, err error)
	MerchantCreateActivity(ctx context.Context, req *core_api.MerchantCreateActivityReq) (resp *core_api.Response, err error)
	MerchantDeleteActivity(ctx context.Context, req *core_api.MerchantDeleteActivityReq) (resp *core_api.Response, err error)
	MerchantUpdateActivity(ctx context.Context, req *core_api.MerchantUpdateActivityReq) (resp *core_api.Response, err error)
	MerchantGetActivity(ctx context.Context, req *core_api.MerchantGetActivityReq) (resp *core_api.MerchantGetActivityResp, err error)
	MerchantTopActivity(ctx context.Context, req *core_api.MerchantTopActivityReq) (resp *core_api.Response, err error)
	MerchantLogin(ctx context.Context, req *core_api.MerchantLoginReq) (resp *core_api.MerchantLoginResp, err error)
	MerchantGetSetting(ctx context.Context, req *core_api.MerchantGetSettingReq) (resp *core_api.MerchantGetSettingResp, err error)
	MerchantUpdateSetting(ctx context.Context, req *core_api.MerchantUpdateSettingReq) (resp *core_api.Response, err error)
	MerchantListBookRecords(ctx context.Context, req *core_api.MerchantListBookRecordsReq) (resp *core_api.MerchantListBookRecordsResp, err error)
	MerchantUpdateInfo(ctx context.Context, req *core_api.MerchantUpdateInfoReq) (resp *core_api.Response, err error)
	MerchantGetInfo(ctx context.Context, req *core_api.MerchantGetInfoReq) (resp *core_api.MerchantGetInfoResp, err error)
	MerchantSetPassword(ctx context.Context, req *core_api.MerchantSetPasswordReq) (resp *core_api.Response, err error)
	GetMerchantInfoByUri(ctx context.Context, c *core_api.GetMerchantInfoByUriReq) (resp *core_api.GetMerchantInfoByUriResp, err error)
	GetAd(ctx context.Context, c *core_api.GetAdReq) (resp *core_api.GetAdResp, err error)
	SetAd(ctx context.Context, c *core_api.SetAdReq) (resp *core_api.Response, err error)
}

type MerchantService struct {
	SystemRpc rpcsystem.IActiManageSystem
	UserRpc   rpcuser.IActiManageUser
	Config    *config.Config
}

var MerchantServiceSet = wire.NewSet(
	wire.Struct(new(MerchantService), "*"),
	wire.Bind(new(IMerchantService), new(*MerchantService)),
)

func (s *MerchantService) MerchantListActivities(ctx context.Context, req *core_api.MerchantListActivitiesReq) (resp *core_api.MerchantListActivitiesResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	listResp, err := s.SystemRpc.ListActivities(ctx, &gensystem.ListActivitiesReq{
		Paging: &genbasic.Paging{
			Page:  req.Paging.Page,
			Limit: req.Paging.Limit,
		},
		Type:       *req.Type,
		MerchantId: userId,
	})
	if err != nil {
		return nil, err
	}
	activities := make([]*core_api.MerchantListActivitiesResp_Item, 0, len(listResp.Activities))
	for _, v := range listResp.Activities {
		activities = append(activities, &core_api.MerchantListActivitiesResp_Item{
			Id:   v.Id,
			Name: v.Name,
			Top:  v.Top,
		})
	}
	resp = &core_api.MerchantListActivitiesResp{
		Code:       0,
		Msg:        "获取成功",
		Activities: activities,
		Total:      listResp.Total,
	}
	return resp, nil
}

func (s *MerchantService) MerchantCreateActivity(ctx context.Context, req *core_api.MerchantCreateActivityReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	rpcReq := &gensystem.CreateActivityReq{
		MerchantId:  userId,
		Name:        req.Name,
		Book:        req.Type,
		Setting:     make([]*gensystem.ActivitySetting, 0, len(req.ActivitySettings)),
		Location:    &gensystem.Location{},
		Top:         req.Top,
		Phone:       req.Phone,
		Description: req.Description,
		Cover:       req.Cover,
	}
	// 活动，有预约时间
	if req.Type == 1 {
		if req.BookStart != nil {
			rpcReq.BookStart = req.GetBookStart()
		}
		if req.BookEnd != nil {
			rpcReq.BookEnd = req.GetBookEnd()
		}
		if req.Notice != nil {
			rpcReq.Notice = req.GetNotice()
		}
	}
	// 活动设置
	err = copier.Copy(&rpcReq.Setting, &req.ActivitySettings)
	if err != nil {
		log.Error("copier.Copy failed", err)
		return nil, consts.ErrInvalidParameter
	}
	// 地点
	err = copier.Copy(&rpcReq.Location, &req.Location)

	response, err := s.SystemRpc.CreateActivity(ctx, rpcReq)
	if err != nil || response.Code != 0 {
		return nil, consts.ErrCreate
	}
	return util.SuccessResp("创建成功")
}

func (s *MerchantService) MerchantDeleteActivity(ctx context.Context, req *core_api.MerchantDeleteActivityReq) (resp *core_api.Response, err error) {
	response, err := s.SystemRpc.DeleteActivity(ctx, &gensystem.DeleteActivityReq{
		Id: req.Id,
	})
	if err != nil || response.Code != 0 {
		return nil, consts.ErrDelete
	}
	return util.SuccessResp("删除成功")
}

func (s *MerchantService) MerchantUpdateActivity(ctx context.Context, req *core_api.MerchantUpdateActivityReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}

	notice := "暂无"
	if req.Notice != nil {
		notice = *req.Notice
	}

	response, err := s.SystemRpc.UpdateActivity(ctx, &gensystem.UpdateActivityReq{
		MerchantId: userId,
		Activity: &gensystem.Activity{
			Id:         req.Id,
			MerchantId: userId,
			Name:       req.Name,
			Cover:      req.Cover,
			Setting: &gensystem.ActivitySetting{
				Max: req.Max,
			},
			Location: &gensystem.Location{
				Text:      req.Location.Text,
				Longitude: req.Location.Longitude,
				Latitude:  req.Location.Latitude,
			},
			Description: req.Description,
			Notice:      notice,
		},
	})
	if err != nil || response.Code != 0 {
		return nil, consts.ErrUpdate
	}
	return util.SuccessResp("更新成功")
}

func (s *MerchantService) MerchantGetActivity(ctx context.Context, req *core_api.MerchantGetActivityReq) (resp *core_api.MerchantGetActivityResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}

	response, err := s.SystemRpc.GetActivity(ctx, &gensystem.GetActivityReq{
		Id: req.Id,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	v := response.Activity
	resp = &core_api.MerchantGetActivityResp{
		Code:     0,
		Msg:      "success",
		Setting:  &core_api.ActivitySetting{},
		Location: &core_api.Location{},
	}
	err = copier.Copy(&resp, &v)
	if err != nil {
		return nil, err
	}
	// 预约设置
	if v.Book == 1 {
		resp.BookStart = v.BookStart
		resp.BookEnd = v.BookEnd
	}

	// 点赞数和浏览数
	fvResp, err := s.UserRpc.GetFavoriteAndViewOfActivity(ctx, &genuser.GetFavoriteAndViewOfActivityReq{
		ActivityId: v.Id,
	})
	if err != nil {
		return nil, err
	}
	resp.Favorite = fvResp.Favorite
	resp.View = fvResp.View

	// 复用预约判断
	bookResp, err := s.UserRpc.CheckBookRecordByUserIdAndActivityId(ctx, &genuser.CheckBookRecordByUserIdAndActivityIdReq{
		UserId:     userId,
		ActivityId: v.Id,
	})
	if err != nil {
		return nil, err
	}
	// 预约总人数
	resp.CurrentBooked = bookResp.CurrentBooked
	return resp, nil
}

func (s *MerchantService) MerchantTopActivity(ctx context.Context, req *core_api.MerchantTopActivityReq) (resp *core_api.Response, err error) {
	response, err := s.SystemRpc.TopActivity(ctx, &gensystem.TopActivityReq{
		Id: req.Id,
	})
	if err != nil || response.Code != 0 {
		return nil, consts.ErrUpdate
	}
	return util.SuccessResp("置顶成功")
}

func (s *MerchantService) MerchantLogin(ctx context.Context, req *core_api.MerchantLoginReq) (resp *core_api.MerchantLoginResp, err error) {
	response, err := s.SystemRpc.MerchantLogin(ctx, &gensystem.MerchantLoginReq{
		AuthId:   req.AuthId,
		Password: req.Password,
	})
	if err != nil || response.Code != 0 {
		log.Error("Login error: ", err)
		return nil, consts.ErrLogin
	}
	token, expire, err := util.GenerateJwtToken(s.Config.Auth.SecretKey, s.Config.Auth.AccessExpire, response.Id)
	if err != nil {
		log.Error("Generate jwt token error: ", err)
		return nil, consts.ErrLogin
	}
	resp = &core_api.MerchantLoginResp{
		Code:        0,
		Msg:         "success",
		Id:          response.Id,
		AccessToken: token,
		ExpireTime:  expire,
	}
	return resp, nil
}

func (s *MerchantService) MerchantGetSetting(ctx context.Context, req *core_api.MerchantGetSettingReq) (resp *core_api.MerchantGetSettingResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.SystemRpc.GetMerchantSetting(ctx, &gensystem.GetMerchantSettingReq{
		MerchantId: userId,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	resp = &core_api.MerchantGetSettingResp{
		Code:   0,
		Msg:    "success",
		Header: &core_api.Header{},
		Footer: &core_api.Footer{},
		Cover:  &core_api.Cover{},
		Id:     response.Setting.Id,
	}
	err = copier.Copy(&resp.Header, &response.Setting.Header)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(&resp.Cover, &response.Setting.Cover)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(&resp.Footer, &response.Setting.Footer)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *MerchantService) MerchantUpdateSetting(ctx context.Context, req *core_api.MerchantUpdateSettingReq) (resp *core_api.Response, err error) {
	rpcReq := &gensystem.UpdateSettingReq{
		Id:     req.Id,
		Header: &gensystem.Header{},
		Cover:  &gensystem.Cover{},
		Footer: &gensystem.Footer{},
	}
	err = copier.Copy(&rpcReq.Header, &req.Header)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(&rpcReq.Cover, &req.Cover)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(&rpcReq.Footer, &req.Footer)
	if err != nil {
		return nil, err
	}
	response, err := s.SystemRpc.UpdateMerchantSetting(ctx, rpcReq)
	if err != nil || response.Code != 0 {
		return nil, err
	}
	return util.SuccessResp("更新成功")
}

func (s *MerchantService) MerchantListBookRecords(ctx context.Context, req *core_api.MerchantListBookRecordsReq) (resp *core_api.MerchantListBookRecordsResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	records, err := s.UserRpc.ListBookRecordsByActivity(ctx, &genuser.ListBookRecordsByActivityReq{
		ActivityId: req.ActivityId,
		Paging: &genbasic.Paging{
			Page:  req.Paging.Page,
			Limit: req.Paging.Limit,
		},
	})
	if err != nil {
		return nil, err
	}
	recordDtos := make([]*core_api.MerchantListBookRecordsResp_BookItem, 0, len(records.Records))
	for _, record := range records.Records {
		var dto core_api.MerchantListBookRecordsResp_BookItem
		err = copier.Copy(&dto, &record)
		if err != nil {
			return nil, err
		}
		// 复制预约人
		rs := make([]*core_api.MerchantListBookRecordsResp_Item, 0, len(record.Reservers))
		for _, reserver := range record.Reservers {
			r := &core_api.MerchantListBookRecordsResp_Item{
				ReserverId: reserver.ReserverId,
				Cancel:     reserver.Cancel,
				Name:       reserver.Name,
				Phone:      reserver.Phone,
			}
			rs = append(rs, r)
		}
		dto.Reservers = rs
		recordDtos = append(recordDtos, &dto)
	}
	return &core_api.MerchantListBookRecordsResp{
		Code:          0,
		Msg:           "success",
		BookRecords:   recordDtos,
		Total:         records.Total,
		CurrentBooked: records.CurrentBooked,
	}, nil
}

func (s *MerchantService) MerchantUpdateInfo(ctx context.Context, req *core_api.MerchantUpdateInfoReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	rpcReq := &gensystem.UpdateMerchantInfoReq{
		Id:          userId,
		Name:        req.Name,
		Logo:        req.Name,
		Description: req.Description,
		Licences:    req.Licences,
		Openings:    make([]*gensystem.Opening, 0, len(req.Openings)),
		Location:    &gensystem.Location{},
	}
	err = copier.Copy(&rpcReq.Openings, &req.Openings)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(&rpcReq.Location, &req.Location)
	if err != nil {
		return nil, err
	}
	response, err := s.SystemRpc.UpdateMerchantInfo(ctx, rpcReq)
	if err != nil || response.Code != 0 {
		return nil, err
	}
	return util.SuccessResp("更新成功")
}

func (s *MerchantService) MerchantGetInfo(ctx context.Context, req *core_api.MerchantGetInfoReq) (resp *core_api.MerchantGetInfoResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.SystemRpc.GetMerchantInfo(ctx, &gensystem.GetMerchantInfoReq{
		Id: userId,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	resp = &core_api.MerchantGetInfoResp{
		Code:     0,
		Msg:      "success",
		Openings: make([]*core_api.Opening, 0, len(response.Openings)),
		Location: &core_api.Location{},
		Uri:      response.Uri,
	}
	err = copier.Copy(&resp, &response)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *MerchantService) MerchantSetPassword(ctx context.Context, req *core_api.MerchantSetPasswordReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.SystemRpc.MerchantSetPassword(ctx, &gensystem.MerchantSetPasswordReq{
		Id:          userId,
		OldPassword: req.OldPassword,
		NewPassword: req.Password,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	return util.SuccessResp("密码修改成功")
}

func (s *MerchantService) GetMerchantInfoByUri(ctx context.Context, req *core_api.GetMerchantInfoByUriReq) (resp *core_api.GetMerchantInfoByUriResp, err error) {
	response, err := s.SystemRpc.GetMerchantInfoByUri(ctx, &gensystem.GetMerchantInfoByUriReq{
		Uri: req.Uri,
	})
	if err != nil {
		return nil, err
	}
	resp = &core_api.GetMerchantInfoByUriResp{
		Code:       0,
		Msg:        "success",
		MerchantId: response.MerchantId,
		Ad:         &core_api.Ad{},
	}
	err = copier.Copy(&resp.Ad, &response.Ad)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *MerchantService) GetAd(ctx context.Context, req *core_api.GetAdReq) (resp *core_api.GetAdResp, err error) {
	response, err := s.SystemRpc.GetAd(ctx, &gensystem.GetAdReq{
		MerchantId: req.MerchantId,
	})
	if err != nil {
		return nil, err
	}
	resp = &core_api.GetAdResp{
		Code: 0,
		Msg:  "success",
		Ad:   &core_api.Ad{},
	}
	err = copier.Copy(&resp.Ad, &response.Ad)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *MerchantService) SetAd(ctx context.Context, req *core_api.SetAdReq) (resp *core_api.Response, err error) {
	rpcReq := &gensystem.SetAdReq{
		MerchantId: req.MerchantId,
		Ad:         &gensystem.Ad{},
	}
	err = copier.Copy(&rpcReq.Ad, &req.Ad)
	if err != nil {
		return nil, err
	}
	response, err := s.SystemRpc.SetAd(ctx, rpcReq)
	if err != nil || response.Code != 0 {
		return nil, err
	}
	return util.SuccessResp(response.Msg)
}
