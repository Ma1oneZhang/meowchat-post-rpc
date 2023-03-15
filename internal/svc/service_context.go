package svc

import (
	"github.com/xh-polaris/meowchat-post-rpc/internal/config"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post"
)

type ServiceContext struct {
	Config    config.Config
	PostModel post.Model
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		PostModel: post.NewModel(c.Mongo.URL, c.Mongo.DB, c.Cache, c.Elasticsearch),
	}
}
