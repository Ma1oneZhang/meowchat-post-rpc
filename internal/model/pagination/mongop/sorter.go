package mongop

import (
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	MongoSorter interface {
		MakeSortOptions(filter bson.M, backward bool) (bson.M, error)
	}

	IdSorter struct {
		ID string `json:"_id"`
	}
)

func (s *IdSorter) MakeSortOptions(filter bson.M, backward bool) (bson.M, error) {
	//构造lastId
	var id primitive.ObjectID
	var err error
	if s == nil {
		if backward {
			id = primitive.NewObjectIDFromTimestamp(time.Unix(0, 0))
		} else {
			id = primitive.NewObjectIDFromTimestamp(time.Unix(math.MaxInt64, 0))
		}
	} else {
		id, err = primitive.ObjectIDFromHex(s.ID)
		if err != nil {
			return nil, err
		}
	}

	sort := bson.M{}
	if backward {
		filter["_id"] = bson.M{"$gt": id}
		sort = bson.M{"_id": 1}
	} else {
		filter["_id"] = bson.M{"$lt": id}
		sort = bson.M{"_id": -1}
	}
	return sort, err
}
