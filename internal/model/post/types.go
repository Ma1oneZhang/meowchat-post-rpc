package post

import (
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal/es"
)

type (
	Post          = internal.Post
	Flag          = internal.PostFlag
	FilterOptions = internal.FilterOptions
)

const (
	OfficialFlag = internal.OfficialFlag
)

const (
	IdSorter    = internal.IdSorter
	ScoreSorter = internal.ScoreSorter
)

var (
	ConvertAllFieldsSearchQuery   = es.ConvertAllFieldsSearchQuery
	ConvertMultiFieldsSearchQuery = es.ConvertMultiFieldsSearchQuery
)
