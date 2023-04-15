package mongo

import (
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"
	"github.com/xh-polaris/paginator-go/mongop"
)

const (
	IdSort = internal.IdSorter
)

var Sorters = map[int64]mongop.MongoSorter{
	IdSort: (*mongop.IdSorter)(nil),
}
