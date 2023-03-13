package paginator

import (
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defaultPageSize = int64(10)
	defaultBackward = false
)

type BasePaginator struct {
	Limit    *int64
	Offset   *int64
	Backward *bool
}

// GenQuery 生成基础查询选项
func (p *BasePaginator) GenQuery() *options.FindOptions {
	if p.Backward == nil {
		p.Backward = &defaultBackward
	}
	if p.Limit == nil {
		p.Limit = &defaultPageSize
	}

	opts := &options.FindOptions{
		Limit: p.Limit,
		Skip:  p.Offset,
	}
	return opts
}

type IdPaginator struct {
	LastId *string
	BasePaginator
}

// GenQuery 生成ID分页查询选项，并将filter在原地更新
func (p *IdPaginator) GenQuery(filter bson.M) (*options.FindOptions, error) {
	opts := p.BasePaginator.GenQuery()
	opts.Sort = bson.M{"_id": -1}

	//构造lastId
	var oid primitive.ObjectID
	var err error
	if p.LastId == nil {
		if *p.Backward {
			oid = primitive.NewObjectIDFromTimestamp(time.Unix(math.MinInt32, 0))
		} else {
			oid = primitive.NewObjectIDFromTimestamp(time.Unix(math.MaxInt32, 0))
		}
	} else {
		oid, err = primitive.ObjectIDFromHex(*p.LastId)
		if err != nil {
			return nil, err
		}
	}

	if *p.Backward {
		filter["_id"] = bson.M{"$gt": oid}
	} else {
		filter["_id"] = bson.M{"$lt": oid}
	}
	return opts, nil
}
