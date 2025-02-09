package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/ActiManage-core-api/biz/application/dto/core_api"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/consts"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/util"
)

type IStsService interface {
	StsApplySignedUrl(ctx context.Context, req *core_api.StsApplySignedUrlReq) (resp *core_api.StsApplySignedUrlResp, err error)
	StsAIModify(ctx context.Context, req *core_api.StsAIModifyReq) (resp *core_api.StsAIModifyResp, err error)
}

type StsService struct {
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
