package provider

import (
	"github.com/google/wire"
	"github.com/xh-polaris/ActiManage-core-api/biz/application/service"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/config"
	rpcsystem "github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/rpc/system"
	rpcuser "github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/rpc/user"
)

var provider *Provider

func Init() {
	var err error
	provider, err = NewProvider()
	if err != nil {
		panic(err)
	}
}

// Provider 提供controller依赖的对象
type Provider struct {
	Config          *config.Config
	MerchantService service.MerchantService
	StsService      service.StsService
	SystemService   service.SystemService
	UserService     service.UserService
}

func Get() *Provider {
	return provider
}

var RPCSet = wire.NewSet(
	rpcsystem.ActiManageSystemSet,
	rpcuser.ActiManageUserSet,
)

var ApplicationSet = wire.NewSet(
	service.MerchantServiceSet,
	service.StsServiceSet,
	service.SystemServiceSet,
	service.UserServiceSet,
)

var DomainSet = wire.NewSet()

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	RPCSet,
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	DomainSet,
	InfrastructureSet,
)
