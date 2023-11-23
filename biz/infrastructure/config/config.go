package config

import (
	"os"

	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Config struct {
	service.ServiceConf
	ListenOn     string
	Cache        cache.CacheConf
	Redis        *redis.RedisConf
	GetFishTimes int64
	Mongo        *struct {
		DB  string
		URL string
	}
	RocketMq *struct {
		URL       []string
		Retry     int
		GroupName string
	}
}

func NewConfig() (*Config, error) {
	c := new(Config)
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "etc/config.yaml"
	}
	err := conf.Load(path, c)
	if err != nil {
		return nil, err
	}
	err = c.SetUp()
	if err != nil {
		return nil, err
	}
	return c, nil
}
