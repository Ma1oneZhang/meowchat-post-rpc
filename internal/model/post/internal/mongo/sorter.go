package mongo

import (
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/pagination/mongop"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"
)

const (
	IdSort = internal.IdSorter
)

var Sorters = map[int64]mongop.MongoSorter{
	IdSort: (*mongop.IdSorter)(nil),
}
