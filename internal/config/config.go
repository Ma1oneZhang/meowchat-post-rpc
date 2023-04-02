package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type ElasticsearchConf struct {
	Addresses []string
	Username  string
	Password  string
}

type Config struct {
	zrpc.RpcServerConf
	Cache cache.CacheConf
	Mongo struct {
		URL string
		DB  string
	}
	Elasticsearch ElasticsearchConf
	RocketMq struct {
		URL []string
		Retry int
		GroupName string
	}
}
