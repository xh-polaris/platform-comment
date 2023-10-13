package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ReplyTo string             `bson:"replyTo,omitempty" json:"replyTo,omitempty"`

	ParentId     string `bson:"parentId,omitempty" json:"parentId,omitempty"`
	FirstLevelId string `bson:"firstLevelId,omitempty" json:"firstLevelId,omitempty"`
	Type         string `bson:"type,omitempty" json:"type,omitempty"`
	Text         string `bson:"text,omitempty" json:"text,omitempty"`
	AuthorId     string `bson:"authorId,omitempty" json:"authorId,omitempty"`

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
