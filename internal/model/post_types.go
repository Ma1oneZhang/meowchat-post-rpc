package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `bson:"title,omitempty"`
	Text     string             `bson:"text,omitempty"`
	CoverUrl string             `bson:"coverUrl,omitempty"`
	Tags     []string           `bson:"tags,omitempty"`
	UserId   string             `bson:"userId,omitempty"`
	UpdateAt time.Time          `bson:"updateAt,omitempty"`
	CreateAt time.Time          `bson:"createAt,omitempty"`
}
