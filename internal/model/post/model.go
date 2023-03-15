package post

import (
	"github.com/xh-polaris/meowchat-post-rpc/internal/config"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal/es"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal/mongo"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type (
	Model interface {
		mongo.PostMongoModel
		es.PostEsModel
	}
	defaultModel struct {
		mongo.PostMongoModel
		es.PostEsModel
	}
)

func NewModel(url, db string, c cache.CacheConf, ec config.ElasticsearchConf) Model {
	return defaultModel{
		PostMongoModel: mongo.NewPostModel(url, db, c),
		PostEsModel:    es.NewPostModel(db, ec, c),
	}
}
