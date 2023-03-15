package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title    string             `bson:"title,omitempty" `
	Text     string             `bson:"text,omitempty"`
	CoverUrl string             `bson:"coverUrl,omitempty"`
	Tags     []string           `bson:"tags,omitempty"`
	UserId   string             `bson:"userId,omitempty"`
	Flags    *PostFlag          `bson:"flags,omitempty"`
	UpdateAt time.Time          `bson:"updateAt,omitempty"`
	CreateAt time.Time          `bson:"createAt,omitempty"`
	// 仅ES查询时使用
	Score_ float64 `bson:"_score,omitempty" json:"_score,omitempty"`
}

const (
	ID       = "_id"
	Title    = "title"
	Text     = "text"
	CoverUrl = "coverUrl"
	Tags     = "tags"
	UserId   = "userId"
	Flags    = "flags"
	UpdateAt = "updateAt"
	CreateAt = "createAt"
)

type PostFlag int64

const (
	OfficialFlag = 1 << 0
)

func (f *PostFlag) SetFlag(flag PostFlag, b bool) *PostFlag {
	if f == nil {
		f = new(PostFlag)
	}
	if b {
		*f |= flag
	} else {
		*f &= ^flag
	}
	return f
}

func (f *PostFlag) GetFlag(flag PostFlag) bool {
	return f != nil && (*f&flag) > 0
}

const (
	IdSorter = iota
	ScoreSorter
)
