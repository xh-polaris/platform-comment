// Code generated by goctl. DO NOT EDIT.
package mapper

import (
	"context"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/data/db"
	"time"

	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var prefixCommentCacheKey = "cache:comment:"

type commentModel interface {
	Insert(ctx context.Context, data *db.Comment) error
	FindOne(ctx context.Context, id string) (*db.Comment, error)
	Update(ctx context.Context, data *db.Comment) error
	Delete(ctx context.Context, id string) error
}

type defaultCommentModel struct {
	conn *monc.Model
}

func newDefaultCommentModel(conn *monc.Model) *defaultCommentModel {
	return &defaultCommentModel{conn: conn}
}

func (m *defaultCommentModel) Insert(ctx context.Context, data *db.Comment) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	key := prefixCommentCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *defaultCommentModel) FindOne(ctx context.Context, id string) (*db.Comment, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidObjectId
	}

	var data db.Comment
	key := prefixCommentCacheKey + id
	err = m.conn.FindOne(ctx, key, &data, bson.M{"_id": oid})
	switch err {
	case nil:
		return &data, nil
	case monc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCommentModel) Update(ctx context.Context, data *db.Comment) error {
	data.UpdateAt = time.Now()
	key := prefixCommentCacheKey + data.ID.Hex()
	_, err := m.conn.ReplaceOne(ctx, key, bson.M{"_id": data.ID}, data)
	return err
}

func (m *defaultCommentModel) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidObjectId
	}
	key := prefixCommentCacheKey + id
	_, err = m.conn.DeleteOne(ctx, key, bson.M{"_id": oid})
	return err
}
