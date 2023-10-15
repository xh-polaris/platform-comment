package mapper

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/config"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/data/db"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

const CommentCollectionName = "comment"

var _ CommentModel = (*customCommentModel)(nil)

type (
	// CommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentModel.
	CommentModel interface {
		commentModel
		FindByAuthorIdAndType(ctx context.Context, authorId string, _type string, skip int64, limit int64) ([]db.Comment, int64, error)
		FindByParent(ctx context.Context, _type string, parentId string, skip int64, limit int64) ([]db.Comment, int64, error)
		FindByFirstLevel(ctx context.Context, firstLevelId string, skip int64, limit int64) ([]db.Comment, int64, error)
		FindByReplyToAndType(ctx context.Context, _type string, replyTo string, skip int64, limit int64) ([]db.Comment, int64, error)
		CountByParent(ctx context.Context, _type string, parentId string) (int64, error)
	}

	customCommentModel struct {
		*defaultCommentModel
	}
)

func (c *customCommentModel) FindByReplyToAndType(ctx context.Context, _type string, replyTo string, skip int64, limit int64) ([]db.Comment, int64, error) {
	var data []db.Comment
	if err := c.conn.Find(ctx, &data,
		bson.M{
			"replyTo": replyTo,
			"type":    _type,
		},
		&mopt.FindOptions{
			Skip:  &skip,
			Limit: &limit,
			Sort:  bson.M{"createAt": -1},
		}); err != nil {
		return nil, 0, err
	}
	total, err := c.conn.CountDocuments(ctx, bson.M{"replyTo": replyTo,
		"type": _type})
	if err != nil {
		return nil, 0, err
	}
	return data, total, err
}

func (c *customCommentModel) FindByParent(ctx context.Context, _type string, parentId string, skip int64, limit int64) ([]db.Comment, int64, error) {
	var data []db.Comment
	if err := c.conn.Find(ctx, &data,
		bson.M{
			"parentId": parentId,
			"type":     _type,
		},
		&mopt.FindOptions{
			Skip:  &skip,
			Limit: &limit,
			Sort:  bson.M{"createAt": -1},
		}); err != nil {
		return nil, 0, err
	}
	total, err := c.conn.CountDocuments(ctx, bson.M{
		"parentId": parentId,
		"type":     _type,
	})
	if err != nil {
		return nil, 0, err
	}
	return data, total, err
}

func (c *customCommentModel) FindByFirstLevel(ctx context.Context, firstLevelId string, skip int64, limit int64) ([]db.Comment, int64, error) {
	var data []db.Comment
	if err := c.conn.Find(ctx, &data,
		bson.M{
			"firstLevelId": firstLevelId,
		},
		&mopt.FindOptions{
			Skip:  &skip,
			Limit: &limit,
			Sort:  bson.M{"createAt": -1},
		}); err != nil {
		return nil, 0, err
	}
	total, err := c.conn.CountDocuments(ctx, bson.M{
		"firstLevelId": firstLevelId,
	})
	if err != nil {
		return nil, 0, err
	}
	return data, total, err
}

func (c *customCommentModel) FindByAuthorIdAndType(ctx context.Context, authorId string, _type string, skip int64, limit int64) ([]db.Comment, int64, error) {
	var data []db.Comment
	if err := c.conn.Find(ctx, &data,
		bson.M{
			"authorId": authorId,
			"type":     _type,
		},
		&mopt.FindOptions{
			Skip:  &skip,
			Limit: &limit,
			Sort:  bson.M{"createAt": -1},
		}); err != nil {
		return nil, 0, err
	}
	total, err := c.conn.CountDocuments(ctx, bson.M{
		"authorId": authorId,
		"type":     _type,
	})
	if err != nil {
		return nil, 0, err
	}
	return data, total, err
}

func (c *customCommentModel) CountByParent(ctx context.Context, _type string, parentId string) (int64, error) {
	total, err := c.conn.CountDocuments(ctx, bson.M{
		"parentId": parentId,
		"type":     _type,
	})
	if err != nil {
		return 0, err
	}
	return total, nil
}

// NewCommentModel returns a model for the mongo.
func NewCommentModel(config *config.Config) CommentModel {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CommentCollectionName, config.Cache)
	return &customCommentModel{
		defaultCommentModel: newDefaultCommentModel(conn),
	}
}

var CommentSet = wire.NewSet(
	NewCommentModel,
)
