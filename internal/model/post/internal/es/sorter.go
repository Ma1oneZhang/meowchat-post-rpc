package es

import (
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"
	"github.com/xh-polaris/paginator-go/esp"
	"math"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	IdSort    = internal.IdSorter
	ScoreSort = internal.ScoreSorter
)

var Sorters = map[int64]esp.EsSorter{
	IdSort:    (*esp.IdSorter)(nil),
	ScoreSort: (*ScoreSorter)(nil),
}

type (
	ScoreSorter struct {
		Score_ float64 `json:"_score"`
		ID     string  `json:"_id"`
	}
)

func (s *ScoreSorter) MakeSortOptions(backward bool) ([]types.SortCombinations, []types.FieldValue, error) {
	var id string
	var score float64
	if s == nil {
		if backward {
			id = primitive.NewObjectIDFromTimestamp(time.Unix(0, 0)).Hex()
			score = 0
		} else {
			id = primitive.NewObjectIDFromTimestamp(time.Unix(math.MaxInt64, 0)).Hex()
			score = math.MaxFloat64
		}
	} else {
		id = s.ID
		score = s.Score_
	}
	var sort []types.SortCombinations
	if !backward {
		sort = append(sort, types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"_score": {Order: &sortorder.Desc},
			},
		})
		sort = append(sort, types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"_id": {Order: &sortorder.Desc},
			},
		})
	} else {
		sort = append(sort, types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"_score": {Order: &sortorder.Asc},
			},
		})
		sort = append(sort, types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"_id": {Order: &sortorder.Asc},
			},
		})
	}
	return sort, []types.FieldValue{score, id}, nil
}
