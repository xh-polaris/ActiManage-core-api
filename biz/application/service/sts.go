package service

import (
	"context"
	"github.com/google/wire"
	gensystem "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/system"
	genuser "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/user"
	"github.com/xh-polaris/ActiManage-core-api/biz/application/dto/core_api"
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
	// TODO 向COS申请url
	return &core_api.StsApplySignedUrlResp{}, nil
}

func (s *StsService) StsAIModify(ctx context.Context, req *core_api.StsAIModifyReq) (resp *core_api.StsAIModifyResp, err error) {
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
