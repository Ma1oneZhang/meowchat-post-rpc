package mongo

import (
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"

	"go.mongodb.org/mongo-driver/bson"
)

// 检验是否完成Filter接口
var _ internal.Filter = (*PostFilter)(nil)

type PostFilter struct {
	m bson.M
	*internal.BaseFilter
}

func MakeBsonFilter(options *internal.FilterOptions) bson.M {
	return (&PostFilter{
		m: bson.M{},
		BaseFilter: &internal.BaseFilter{
			FilterOptions: options,
		},
	}).toBson()
}

func (f *PostFilter) toBson() bson.M {
	f.CheckOnlyUserId()
	f.CheckOnlyOfficial()
	f.CheckFlags()
	return f.m
}

func (f *PostFilter) CheckFlags() {
	if f.MustFlags != nil {
		f.m[internal.Flags] = bson.M{"$bitsAllSet": *f.MustFlags}
	}
	if f.MustNotFlags != nil {
		or, exist := f.m["$or"]
		if !exist {
			or = bson.A{}
		}

		_ = append(or.(bson.A), bson.M{
			internal.Flags: bson.M{
				"$bitsAllClear": *f.MustNotFlags},
		}, bson.M{
			internal.Flags: bson.M{
				"$exists": false,
			},
		})
		f.m["$or"] = or
	}
}

func (f *PostFilter) CheckOnlyUserId() {
	if f.OnlyUserId != nil {
		f.m[internal.UserId] = *f.OnlyUserId
	}
}
