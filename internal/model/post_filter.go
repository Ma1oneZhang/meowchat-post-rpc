package model

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.mongodb.org/mongo-driver/bson"
)

type Filter struct {
	mustFlags    *PostFlag
	mustNotFlags *PostFlag
	UserId       *string
	OnlyOfficial *bool
}

func (f *Filter) checkOnlyOfficial() {
	if f.OnlyOfficial != nil {
		f.mustFlags = f.mustFlags.SetOfficial(*f.OnlyOfficial)
	}
}

// EsFilter is filter used for Elasticsearch
type EsFilter struct {
	q []types.Query
	Filter
}

func (f *EsFilter) toEsQuery() []types.Query {
	f.q = make([]types.Query, 0)
	f.checkUserId()
	f.checkOnlyOfficial()
	f.checkFlags()
	return f.q
}

func (f *EsFilter) checkFlags() {
	if f.mustFlags != nil {
		f.q = append(f.q, types.Query{
			Script: &types.ScriptQuery{
				Script: types.InlineScript{
					Source: fmt.Sprintf("doc['%s'].size() != 0 && "+
						"(doc['%s'].value & params.%s) == params.%s", Flags, Flags, Flags, Flags),
					Params: map[string]any{
						Flags: *f.mustFlags,
					},
				},
			},
		})
	}
	if f.mustNotFlags != nil {
		f.q = append(f.q, types.Query{
			Script: &types.ScriptQuery{
				Script: types.InlineScript{
					Source: fmt.Sprintf("doc['%s'].size() == 0 || "+
						"(doc['%s'].value & params.%s) == 0", Flags, Flags, Flags),
					Params: map[string]any{
						Flags: *f.mustFlags,
					},
				},
			},
		})
	}
}

func (f *EsFilter) checkUserId() {
	if f.UserId != nil {
		f.q = append(f.q, types.Query{
			Term: map[string]types.TermQuery{
				UserId: {Value: *f.UserId},
			},
		})
	}
}

type MongoFilter struct {
	m bson.M
	Filter
}

func (f *MongoFilter) toBson() bson.M {
	f.m = bson.M{}
	f.checkUserId()
	f.checkOnlyOfficial()
	f.checkFlags()
	return f.m
}

func (f *MongoFilter) checkFlags() {
	if f.mustFlags != nil {
		f.m[Flags] = bson.M{"$bitsAllSet": *f.mustFlags}
	}
	if f.mustNotFlags != nil {
		or, exist := f.m["$or"]
		if !exist {
			or = bson.A{}
		}

		arr := or.(bson.A)
		arr = append(arr, bson.M{
			Flags: bson.M{
				"$bitsAllClear": *f.mustNotFlags},
		}, bson.M{
			Flags: bson.M{
				"$exists": false,
			},
		})
		f.m["$or"] = or
	}
}

func (f *MongoFilter) checkUserId() {
	if f.UserId != nil {
		f.m[UserId] = *f.UserId
	}
}
