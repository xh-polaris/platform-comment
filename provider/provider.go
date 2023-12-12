package provider

import (
	"github.com/google/wire"

	"github.com/xh-polaris/platform-comment/biz/application/service"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/config"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/mapper"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/stores/redis"
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	InfrastructureSet,
)

var ApplicationSet = wire.NewSet(
	service.CommentSet,
)

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	redis.NewRedis,
	MapperSet,
)

var MapperSet = wire.NewSet(
	mapper.CommentSet,
	mapper.HistorySet,
)
