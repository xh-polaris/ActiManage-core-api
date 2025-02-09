package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/ActiManage-core-api/biz/application/dto/core_api"
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
	//TODO implement me
	panic("implement me")
}

func (s *StsService) StsAIModify(ctx context.Context, req *core_api.StsAIModifyReq) (resp *core_api.StsAIModifyResp, err error) {
	//TODO implement me
	panic("implement me")
}
