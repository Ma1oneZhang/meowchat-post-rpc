package es

import (
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"
	"github.com/xh-polaris/paginator-go/esp"
)

const (
	ScoreSort = internal.ScoreSorter
)

var Sorters = map[int64]esp.EsSorter{
	ScoreSort: (*esp.ScoreSorter)(nil),
}
