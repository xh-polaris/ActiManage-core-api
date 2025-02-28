package service

import (
	"context"
	"fmt"
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
	"time"
)

type IUserService interface {
	Login(ctx context.Context, req *core_api.LoginReq) (resp *core_api.LoginResp, err error)
	SignUp(ctx context.Context, req *core_api.SignUpReq) (resp *core_api.SignUpResp, err error)
	GetSetting(ctx context.Context, req *core_api.GetSettingReq) (resp *core_api.GetSettingResp, err error)
	ListActivities(ctx context.Context, req *core_api.ListActivitiesReq) (resp *core_api.ListActivitiesResp, err error)
	GetActivity(ctx context.Context, req *core_api.GetActivityReq) (resp *core_api.GetActivityResp, err error)
	DoFavorite(ctx context.Context, req *core_api.DoFavoriteReq) (resp *core_api.Response, err error)
	CancelFavorite(ctx context.Context, req *core_api.CancelFavoriteReq) (resp *core_api.Response, err error)
	CreateBooking(ctx context.Context, req *core_api.CreateBookingReq) (resp *core_api.Response, err error)
	CancelBookRecord(ctx context.Context, req *core_api.CancelBookRecordReq) (resp *core_api.Response, err error)
	ListActivitiesByBookRecords(ctx context.Context, req *core_api.ListActivitiesByBookRecordsReq) (resp *core_api.ListActivitiesByBookRecordsResp, err error)
	ListReservers(ctx context.Context, req *core_api.ListReserversReq) (resp *core_api.ListReserversResp, err error)
	CreateReserver(ctx context.Context, req *core_api.CreateReserverReq) (resp *core_api.Response, err error)
	DeleteReserver(ctx context.Context, req *core_api.DeleteReserverReq) (resp *core_api.Response, err error)
	UpdateReserver(ctx context.Context, req *core_api.UpdateReserverReq) (resp *core_api.Response, err error)
	GetUserInfo(ctx context.Context, req *core_api.GetUserInfoReq) (resp *core_api.GetUserInfoResp, err error)
	UpdateUserInfo(ctx context.Context, req *core_api.UpdateUserInfoReq) (resp *core_api.Response, err error)
	UpdateNotice(ctx context.Context, req *core_api.UpdateNoticeReq) (resp *core_api.Response, err error)
	GetMerchantInfo(ctx context.Context, req *core_api.GetMerchantInfoReq) (resp *core_api.GetMerchantInfoResp, err error)
	ListActivitiesByView(ctx context.Context, req *core_api.ListActivitiesByViewReq) (resp *core_api.ListActivitiesByViewResp, err error)
}

type UserService struct {
	UserRpc   rpcuser.IActiManageUser
	SystemRpc rpcsystem.IActiManageSystem
	Config    *config.Config
}

var UserServiceSet = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),
)

func (s UserService) Login(ctx context.Context, req *core_api.LoginReq) (resp *core_api.LoginResp, err error) {
	// 参数校验,authId不能为空
	if req.AuthId == "" {
		return nil, consts.ErrSignUp
	}

	// 这里的逻辑其实有点乱，先验证验证码再验证密码就好，懒得改了
	// 校验验证码
	verifyCheck := "false"
	if req.VerifyCode != nil {
		checkResp, err := s.SystemRpc.StsCheckVerifyCode(ctx, &gensystem.StsCheckVerifyCodeReq{
			VerifyId:   "verify:" + req.MerchantId + ":" + req.AuthId,
			VerifyCode: *req.VerifyCode,
		})
		if err != nil || checkResp.Code != 0 {
			return nil, consts.ErrVerifyCode
		}
		verifyCheck = "true"
	}

	rpcReq := &genuser.UserLoginReq{
		AuthId:     req.AuthId,
		AuthType:   req.AuthType,
		VerifyCode: &verifyCheck,
		Password:   req.Password,
		MerchantId: req.MerchantId,
	}
	if verifyCheck == "false" && req.Password != nil {
		rpcReq.VerifyCode = nil
	}
	response, err := s.UserRpc.UserLogin(ctx, rpcReq)

	if err != nil || response.Id == "" {
		log.Error("Login error: ", err)
		return nil, consts.ErrLogin
	}
	token, expire, err := util.GenerateJwtToken(s.Config.Auth.SecretKey, s.Config.Auth.AccessExpire, response.Id)
	if err != nil {
		log.Error("Generate jwt token error: ", err)
		return nil, consts.ErrLogin
	}
	resp = &core_api.LoginResp{
		Code:        0,
		Msg:         "success",
		Id:          response.Id,
		AccessToken: token,
		ExpireTime:  expire,
	}
	return resp, nil
}

func (s UserService) SignUp(ctx context.Context, req *core_api.SignUpReq) (resp *core_api.SignUpResp, err error) {
	// 参数校验,authId不能为空
	if req.AuthId == "" {
		return nil, consts.ErrSignUp
	}

	// 校验验证码
	verifyCheck := "false"
	checkResp, err := s.SystemRpc.StsCheckVerifyCode(ctx, &gensystem.StsCheckVerifyCodeReq{
		VerifyId:   "verify:" + req.MerchantId + ":" + req.AuthId,
		VerifyCode: req.VerifyCode,
	})
	if err != nil || checkResp.Code != 0 {
		return nil, consts.ErrVerifyCode
	}
	verifyCheck = "true"

	response, err := s.UserRpc.UserSignUp(ctx, &genuser.UserSignUpReq{
		MerchantId: req.MerchantId,
		Name:       &req.Name,
		AuthId:     req.AuthId,
		AuthType:   req.AuthType,
		VerifyCode: verifyCheck,
		Password:   req.Password,
		Gender:     0,
	})
	if err != nil || response.Code != 0 || response.Id == "" {
		return nil, consts.ErrSignUp
	}
	token, expire, err := util.GenerateJwtToken(s.Config.Auth.SecretKey, s.Config.Auth.AccessExpire, response.Id)
	if err != nil {
		log.Error("Generate jwt token error: ", err)
		return nil, consts.ErrLogin
	}
	resp = &core_api.SignUpResp{
		Code:        0,
		Msg:         "success",
		Id:          response.Id,
		AccessToken: token,
		ExpireTime:  expire,
	}
	return resp, nil
}

func (s UserService) GetSetting(ctx context.Context, req *core_api.GetSettingReq) (resp *core_api.GetSettingResp, err error) {
	response, err := s.SystemRpc.GetMerchantSetting(ctx, &gensystem.GetMerchantSettingReq{MerchantId: req.MerchantId})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	resp = &core_api.GetSettingResp{
		Code:   0,
		Msg:    "success",
		Header: &core_api.Header{},
		Footer: &core_api.Footer{},
		Cover:  &core_api.Cover{},
	}
	err = copier.Copy(&resp, &response.Setting)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s UserService) ListActivities(ctx context.Context, req *core_api.ListActivitiesReq) (resp *core_api.ListActivitiesResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}

	rpcReq := &gensystem.ListActivitiesReq{
		Paging: &genbasic.Paging{
			Page:  req.Paging.Page,
			Limit: req.Paging.Limit,
		},
		MerchantId: req.MerchantId,
	}
	if req.Type != nil {
		rpcReq.Type = *req.Type
	}
	response, err := s.SystemRpc.ListActivities(ctx, rpcReq)
	if err != nil || response.Code != 0 {
		return nil, err
	}
	activities := make([]*core_api.ListActivitiesResp_Item, 0)
	for _, v := range response.Activities {
		activity := &core_api.ListActivitiesResp_Item{
			Setting:  &core_api.ActivitySetting{},
			Location: &core_api.Location{},
		}
		// 活动设置
		err = copier.Copy(&activity, &v)
		if err != nil {
			return nil, err
		}
		// 预约设置
		if v.Book == 1 {
			activity.BookStart = v.BookStart
			activity.BookEnd = v.BookEnd
		}

		// 点赞数和浏览数
		fvResp, err := s.UserRpc.GetFavoriteAndViewOfActivity(ctx, &genuser.GetFavoriteAndViewOfActivityReq{
			ActivityId: v.Id,
		})
		if err != nil {
			return nil, err
		}
		activity.Favorite = fvResp.Favorite
		activity.View = fvResp.View

		// 是否收藏
		favoriteResp, err := s.UserRpc.CheckFavorite(ctx, &genuser.CheckFavoriteReq{
			UserId:     userId,
			ActivityId: v.Id,
		})
		if err == nil && favoriteResp.Code == 0 {
			activity.IsFavorite = 1
		}

		activities = append(activities, activity)
	}
	resp = &core_api.ListActivitiesResp{
		Code:       0,
		Msg:        "查询成功",
		Activities: activities,
		Total:      response.Total,
	}
	return resp, nil
}

func (s UserService) GetActivity(ctx context.Context, req *core_api.GetActivityReq) (resp *core_api.GetActivityResp, err error) {
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
	resp = &core_api.GetActivityResp{
		Code:        0,
		Msg:         "success",
		Id:          v.Id,
		MerchantId:  v.MerchantId,
		Name:        v.Name,
		Cover:       v.Cover,
		Book:        v.Book,
		Setting:     &core_api.ActivitySetting{},
		Location:    &core_api.Location{},
		Top:         v.Top,
		Phone:       v.Phone,
		Description: v.Description,
		Notice:      v.Notice,
		CreateTime:  v.CreateTime,
		UpdateTime:  v.UpdateTime,
		Status:      v.Status,
	}
	err = copier.Copy(&resp.Setting, &v.Setting)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(&resp.Location, &v.Location)
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

	// 预约判断
	bookResp, err := s.UserRpc.CheckBookRecordByUserIdAndActivityId(ctx, &genuser.CheckBookRecordByUserIdAndActivityIdReq{
		UserId:     userId,
		ActivityId: v.Id,
	})
	if err != nil {
		return nil, err
	}
	resp.Booked = bookResp.Booked
	// 预约总人数
	resp.CurrentBooked = bookResp.CurrentBooked

	// 是否收藏
	favoriteResp, err := s.UserRpc.CheckFavorite(ctx, &genuser.CheckFavoriteReq{
		UserId:     userId,
		ActivityId: v.Id,
	})
	if err == nil && favoriteResp.Code == 0 {
		resp.IsFavorite = 1
	}
	return resp, nil
}

func (s UserService) DoFavorite(ctx context.Context, req *core_api.DoFavoriteReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.UserRpc.DoFavorite(ctx, &genuser.DoFavoriteReq{
		UserId:     userId,
		ActivityId: req.Id,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	return util.SuccessResp("收藏成功")
}

func (s UserService) CancelFavorite(ctx context.Context, req *core_api.CancelFavoriteReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.UserRpc.CancelFavorite(ctx, &genuser.CancelFavoriteReq{
		UserId:     userId,
		ActivityId: req.Id,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	return util.SuccessResp("取消成功")
}

func (s UserService) CreateBooking(ctx context.Context, req *core_api.CreateBookingReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}

	now := time.Now().Unix()

	// 获取活动
	actiResp, err := s.SystemRpc.GetActivity(ctx, &gensystem.GetActivityReq{
		Id: req.ActivityId,
	})
	if err != nil || actiResp.Code != 0 {
		return nil, err
	}
	activity := actiResp.Activity

	check, err := s.UserRpc.CheckBookRecordByUserIdAndActivityId(ctx, &genuser.CheckBookRecordByUserIdAndActivityIdReq{
		UserId:     userId,
		ActivityId: activity.Id,
	})
	if err != nil {
		return nil, err
	}

	// 校验活动能否预约
	switch {
	case activity.Book == 0:
		fallthrough
	// 预约时间判断
	case activity.BookStart > now || now > activity.BookEnd:
		return nil, consts.ErrExpireBook
	// 重复预约
	case check.Booked == 1:
		return nil, consts.ErrRepeatBook
	// 预约人数判断
	case check.CurrentBooked+int64(len(req.ReserverIds)) > activity.Setting.Max:
		return nil, consts.ErrFullBook
	}

	// 创建预约
	bookResp, err := s.UserRpc.CreateBookRecord(ctx, &genuser.CreateBookRecordReq{
		UserId:      userId,
		ActivityId:  req.ActivityId,
		ReserverIds: req.ReserverIds,
		Arrival:     req.Arrival,
		Remark:      req.Remark,
		MerchantId:  req.MerchantId,
	})
	if err != nil || bookResp.Code != 0 {
		return nil, err
	}

	// 创建回执
	receiptResp, err := s.UserRpc.CreateReceipt(ctx, &genuser.CreateReceiptReq{
		UserId:     userId,
		ActivityId: req.ActivityId,
		Msg:        fmt.Sprintf("活动预约成功!\n\r活动时间为%s-%s,活动地点为%s,预约到店时间为%s,请按时到达", time.Unix(activity.Setting.Start, 0).String(), time.Unix(activity.Setting.End, 0).String(), activity.Location.Text, time.Unix(req.Arrival, 0).String()),
	})
	if err != nil || receiptResp.Code != 0 {
		log.Error("回执创建失败:", err)
	}

	return util.SuccessResp("预约成功")
}

func (s UserService) CancelBookRecord(ctx context.Context, req *core_api.CancelBookRecordReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.UserRpc.CancelBookRecord(ctx, &genuser.CancelBookRecordReq{
		BookRecordId: req.BookRecordId,
		ReserverId:   req.ReserverId,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	return util.SuccessResp("取消成功")
}

func (s UserService) ListActivitiesByBookRecords(ctx context.Context, req *core_api.ListActivitiesByBookRecordsReq) (resp *core_api.ListActivitiesByBookRecordsResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	bookResp, err := s.UserRpc.ListBookRecordsByUser(ctx, &genuser.ListBookRecordsByUserReq{
		UserId: userId,
		Paging: &genbasic.Paging{
			Page:  req.Paging.Page,
			Limit: req.Paging.Limit,
		},
		Type: req.Type,
	})
	if err != nil || bookResp.Code != 0 {
		return nil, err
	}
	activities := make([]*core_api.ListActivitiesByBookRecordsResp_Item, 0, len(bookResp.Records))
	for _, r := range bookResp.Records {
		response, err := s.SystemRpc.GetActivity(ctx, &gensystem.GetActivityReq{
			Id: r.ActivityId,
		})
		if err != nil || response.Code != 0 {
			return nil, err
		}
		v := response.Activity
		activity := &core_api.ListActivitiesByBookRecordsResp_Item{
			Setting:  &core_api.ActivitySetting{},
			Location: &core_api.Location{},
		}
		err = copier.Copy(&activity, &v)
		if err != nil {
			return nil, err
		}
		// 预约设置
		if v.Book == 1 {
			activity.BookStart = v.BookStart
			activity.BookEnd = v.BookEnd
		}

		// 点赞数和浏览数
		fvResp, err := s.UserRpc.GetFavoriteAndViewOfActivity(ctx, &genuser.GetFavoriteAndViewOfActivityReq{
			ActivityId: v.Id,
		})
		if err != nil {
			return nil, err
		}
		activity.Favorite = fvResp.Favorite
		activity.View = fvResp.View
		activities = append(activities, activity)
	}
	return &core_api.ListActivitiesByBookRecordsResp{
		Code:       0,
		Msg:        "success",
		Activities: activities,
		Total:      bookResp.Total,
	}, nil
}

func (s UserService) ListReservers(ctx context.Context, req *core_api.ListReserversReq) (resp *core_api.ListReserversResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.UserRpc.ListReservers(ctx, &genuser.ListReserversReq{
		UserId: userId,
		Paging: &genbasic.Paging{
			Page:  req.Paging.Page,
			Limit: req.Paging.Limit,
		},
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	reservers := make([]*core_api.Reserver, 0)
	for _, v := range response.Reservers {
		reserver := &core_api.Reserver{
			Id:         v.Id,
			UserId:     v.UserId,
			Name:       v.Name,
			Gender:     v.Gender,
			Relation:   v.Relation,
			Phone:      v.Phone,
			Email:      v.Email,
			Birth:      v.Birth,
			CreateTime: v.CreateTime,
			UpdateTime: v.UpdateTime,
			Status:     v.Status,
		}
		reservers = append(reservers, reserver)
	}
	return &core_api.ListReserversResp{
		Code:      0,
		Msg:       "success",
		Reservers: reservers,
		Total:     response.Total,
	}, nil
}

func (s UserService) CreateReserver(ctx context.Context, req *core_api.CreateReserverReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.UserRpc.CreateReserver(ctx, &genuser.CreateReserverReq{
		UserId:   userId,
		Name:     req.Name,
		Gender:   req.Gender,
		Relation: req.Relation,
		Phone:    req.Phone,
		Email:    req.Email,
		Birth:    req.Birth,
	})
	if err != nil || response.Code != 0 {
		return nil, consts.ErrCreate
	}
	return util.SuccessResp("创建成功")
}

func (s UserService) DeleteReserver(ctx context.Context, req *core_api.DeleteReserverReq) (resp *core_api.Response, err error) {
	response, err := s.UserRpc.DeleteReserver(ctx, &genuser.DeleteReserverReq{
		ReserverId: req.Id,
	})
	if err != nil || response.Code != 0 {
		return nil, consts.ErrDelete
	}
	return util.SuccessResp("删除成功")
}

func (s UserService) UpdateReserver(ctx context.Context, req *core_api.UpdateReserverReq) (resp *core_api.Response, err error) {
	reserver := &genuser.Reserver{}
	err = copier.Copy(reserver, &req.Reserver)
	if err != nil {
		return nil, err
	}
	response, err := s.UserRpc.UpdateReserver(ctx, &genuser.UpdateReserverReq{
		Reserver: reserver,
	})
	if err != nil || response.Code != 0 {
		return nil, consts.ErrUpdate
	}
	return util.SuccessResp("更新成功")
}

func (s UserService) GetUserInfo(ctx context.Context, req *core_api.GetUserInfoReq) (resp *core_api.GetUserInfoResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.UserRpc.GetUserInfo(ctx, &genuser.GetUserInfoReq{
		UserId: userId,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	resp = &core_api.GetUserInfoResp{}
	err = copier.Copy(&resp, &response)
	if err != nil {
		log.Error("copier.Copy failed: ", err)
		return nil, err
	}
	resp.Code = 0
	resp.Msg = "success"
	return resp, nil
}

func (s UserService) UpdateUserInfo(ctx context.Context, req *core_api.UpdateUserInfoReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	rpcReq := &genuser.UpdateUserInfoReq{
		Id:          userId,
		Name:        req.Name,
		Gender:      req.Gender,
		Birth:       nil,
		Description: req.Description,
		Avatar:      req.Avatar,
	}
	if req.Birth != nil {
		rpcReq.Birth = req.Birth
	}
	response, err := s.UserRpc.UpdateUserInfo(ctx, rpcReq)
	if err != nil || response.Code != 0 {
		return nil, consts.ErrUpdate
	}
	return util.SuccessResp("更新成功")
}

func (s UserService) UpdateNotice(ctx context.Context, req *core_api.UpdateNoticeReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.UserRpc.SetNotice(ctx, &genuser.SetNoticeReq{
		Id: userId,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	return util.SuccessResp("success")
}

func (s UserService) GetMerchantInfo(ctx context.Context, req *core_api.GetMerchantInfoReq) (resp *core_api.GetMerchantInfoResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.SystemRpc.GetMerchantInfo(ctx, &gensystem.GetMerchantInfoReq{
		Id: req.MerchantId,
	})
	if err != nil || response.Code != 0 {
		return nil, err
	}
	resp = &core_api.GetMerchantInfoResp{
		Code:     0,
		Msg:      "success",
		Openings: make([]*core_api.Opening, 0, len(response.Openings)),
		Location: &core_api.Location{},
	}
	err = copier.Copy(&resp, &response)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s UserService) ListActivitiesByView(ctx context.Context, req *core_api.ListActivitiesByViewReq) (resp *core_api.ListActivitiesByViewResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	viewResp, err := s.UserRpc.ListActivityIdsByView(ctx, &genuser.ListActivityIdsByViewReq{
		UserId: userId,
		Paging: &genbasic.Paging{
			Page:  req.Paging.Page,
			Limit: req.Paging.Limit,
		},
	})
	if err != nil {
		return nil, err
	}

	activities := make([]*core_api.ListActivitiesByViewResp_Item, 0, len(viewResp.Ids))
	for _, id := range viewResp.Ids {
		response, err := s.SystemRpc.GetActivity(ctx, &gensystem.GetActivityReq{
			Id: id,
		})
		if err != nil || response.Code != 0 {
			return nil, err
		}
		v := response.Activity
		activity := &core_api.ListActivitiesByViewResp_Item{
			Setting:  &core_api.ActivitySetting{},
			Location: &core_api.Location{},
		}
		err = copier.Copy(&activity, &v)
		if err != nil {
			return nil, err
		}
		// 预约设置
		if v.Book == 1 {
			activity.BookStart = v.BookStart
			activity.BookEnd = v.BookEnd
		}

		// 点赞数和浏览数
		fvResp, err := s.UserRpc.GetFavoriteAndViewOfActivity(ctx, &genuser.GetFavoriteAndViewOfActivityReq{
			ActivityId: v.Id,
		})
		if err != nil {
			return nil, err
		}
		activity.Favorite = fvResp.Favorite
		activity.View = fvResp.View
		activities = append(activities, activity)
	}
	return &core_api.ListActivitiesByViewResp{
		Code:       0,
		Msg:        "success",
		Activities: activities,
		Total:      viewResp.Total,
	}, nil
}
