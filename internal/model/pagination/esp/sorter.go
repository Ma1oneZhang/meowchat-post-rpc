package esp

import (
	"math"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	EsSorter interface {
		MakeSortOptions(backward bool) ([]types.SortCombinations, []types.FieldValue, error)
	}

	IdSorter struct {
		ID string `json:"_id"`
	}
)

func (s *IdSorter) MakeSortOptions(backward bool) ([]types.SortCombinations, []types.FieldValue, error) {
	var id string
	if s == nil {
		if backward {
			id = primitive.NewObjectIDFromTimestamp(time.Unix(0, 0)).Hex()
		} else {
			id = primitive.NewObjectIDFromTimestamp(time.Unix(math.MaxInt64, 0)).Hex()
		}
	} else {
		id = s.ID
	}
	var sort []types.SortCombinations
	if !backward {
		sort = append(sort, types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"_id": {Order: &sortorder.Desc},
			},
		})
	} else {
		sort = append(sort, types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"_id": {Order: &sortorder.Asc},
			},
		})
	}
	return sort, []types.FieldValue{id}, nil
}
