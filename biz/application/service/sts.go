package service

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/google/wire"
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

type IStsService interface {
	StsApplySignedUrl(ctx context.Context, req *core_api.StsApplySignedUrlReq) (resp *core_api.StsApplySignedUrlResp, err error)
	StsAIModify(ctx context.Context, req *core_api.StsAIModifyReq) (resp *core_api.StsAIModifyResp, err error)
}

type StsService struct {
	SystemRpc rpcsystem.IActiManageSystem
	UserRpc   rpcuser.IActiManageUser
}

var StsServiceSet = wire.NewSet(
	wire.Struct(new(StsService), "*"),
	wire.Bind(new(IStsService), new(*StsService)),
)

func (s *StsService) StsApplySignedUrl(ctx context.Context, req *core_api.StsApplySignedUrlReq) (resp *core_api.StsApplySignedUrlResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	clientOptions := []oss.ClientOption{oss.Region("cn-hongkong")}
	client, err := oss.New("https://oss-cn-hongkong.aliyuncs.com", config.GetConfig().OSSAccessKeyID, config.GetConfig().OSSAccessKeySecret, clientOptions...)
	if err != nil {
		log.Error("oss client init fail", err)
		return nil, err
	}
	bucketName := "actimanage-statics"
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Error("oss client get bucket fail", err)
		return nil, err
	}
	if req.Prefix != "" {
		req.Prefix += "/"
	}
	path := "actimanage/" + userId + "/" + req.Prefix + uuid.New().String() + req.GetSuffix()
	// 有效期30s
	signedURL, err := bucket.SignURL(path, oss.HTTPPut, 60)
	if err != nil {
		log.Error("oss client sign url fail", err)
		return nil, err
	}
	return &core_api.StsApplySignedUrlResp{
		Code: 0,
		Msg:  "success",
		Url:  signedURL,
	}, nil
}

func (s *StsService) StsAIModify(ctx context.Context, req *core_api.StsAIModifyReq) (resp *core_api.StsAIModifyResp, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	httpClient := util.NewHttpClient()
	response, err := httpClient.CallGLM(req.Text, req.Lang)
	if err != nil {
		return nil, consts.ErrCall
	}
	message := response["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})
	text, ok := message["content"].(string)
	if !ok {
		return nil, consts.ErrCall
	}
	return &core_api.StsAIModifyResp{
		Code: 0,
		Msg:  "模型调用成功",
		Text: text,
	}, nil
}

func (s *StsService) StsSendVerifyCode(ctx context.Context, req *core_api.StsSendVerifyCodeReq) (resp *core_api.Response, err error) {
	response, err := s.SystemRpc.StsSendVerifyCode(ctx, &gensystem.StsSendVerifyCodeReq{
		AuthId:   req.AuthId,
		AuthType: req.AuthType,
		Purpose:  req.Purpose,
	})
	if err != nil || response.Code != 0 {
		log.Error("验证码发送失败:", err)
		return nil, consts.ErrSend
	}
	return util.SuccessResp("发送成功")
}

func (s *StsService) StsView(ctx context.Context, req *core_api.StsViewReq) (resp *core_api.Response, err error) {
	userId, err := adaptor.ExtractUserId(ctx)
	if err != nil || userId == "" {
		return nil, consts.ErrNotAuthentication
	}
	response, err := s.UserRpc.CreateView(ctx, &genuser.CreateViewReq{
		TargetId: req.TargetId,
		Type:     req.Type,
	})
	if err != nil || response.Code != 0 {
		log.Error("记录访问失败:", err)
		return nil, consts.ErrSend
	}
	return util.SuccessResp("记录成功")
}
