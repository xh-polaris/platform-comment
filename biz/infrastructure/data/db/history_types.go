package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type History struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type     string             `bson:"type,omitempty" json:"type,omitempty"`
	Data     Comment            `bson:"data,omitempty" json:"data,omitempty"`
	UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
