package system

import (
	"github.com/google/wire"
	system "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/system/systemservice"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/config"
	"github.com/xh-polaris/gopkg/kitex/client"
)

type IActiManageSystem interface {
	system.Client
}

type ActiManageSystem struct {
	system.Client
}

var ActiManageSystemSet = wire.NewSet(
	NewActiManageSystem,
	wire.Struct(new(ActiManageSystem), "*"),
	wire.Bind(new(IActiManageSystem), new(*ActiManageSystem)),
)

func NewActiManageSystem(config *config.Config) system.Client {
	return client.NewClient(config.Name, "ActiManage.system", system.NewClient)
}
