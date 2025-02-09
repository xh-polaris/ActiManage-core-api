package openapi_charge

import (
	"github.com/google/wire"
	user "github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/user/userservice"
	"github.com/xh-polaris/ActiManage-core-api/biz/infrastructure/config"
	"github.com/xh-polaris/gopkg/kitex/client"
)

type IActiManageUser interface {
	user.Client
}

type ActiManageUser struct {
	user.Client
}

var ActiManageUserSet = wire.NewSet(
	NewActiManageUser,
	wire.Struct(new(ActiManageUser), "*"),
	wire.Bind(new(IActiManageUser), new(*ActiManageUser)),
)

func NewActiManageUser(config *config.Config) user.Client {
	return client.NewClient(config.Name, "ActiManage.user", user.NewClient)
}
