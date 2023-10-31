package service

import (
	"context"
	"strconv"
	"time"

	"github.com/google/wire"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/platform/comment"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/xh-polaris/platform-comment/biz/infrastructure/config"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/consts"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/data/db"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/mapper"
)

type ICommentService interface {
	CreateComment(ctx context.Context, req *comment.CreateCommentReq) (resp *comment.CreateCommentResp, err error)
	UpdateComment(ctx context.Context, req *comment.UpdateCommentReq) (resp *comment.UpdateCommentResp, err error)
	DeleteComment(ctx context.Context, req *comment.DeleteCommentByIdReq) (resp *comment.DeleteCommentByIdResp, err error)
	ListCommentByParent(ctx context.Context, req *comment.ListCommentByParentReq) (resp *comment.ListCommentByParentResp, err error)
	CountCommentByParent(ctx context.Context, req *comment.CountCommentByParentReq) (resp *comment.CountCommentByParentResp, err error)
	RetrieveCommentById(ctx context.Context, req *comment.RetrieveCommentByIdReq) (resp *comment.RetrieveCommentByIdResp, err error)
	ListCommentByAuthorIdAndType(ctx context.Context, req *comment.ListCommentByAuthorIdAndTypeReq) (resp *comment.ListCommentByAuthorIdAndTypeResp, err error)
	ListCommentByReplyToAndType(ctx context.Context, req *comment.ListCommentByReplyToAndTypeReq) (resp *comment.ListCommentByReplyToAndTypeResp, err error)
}

type CommentService struct {
	Config       *config.Config
	CommentModel mapper.CommentModel
	HistoryModel mapper.HistoryModel
	Redis        *redis.Redis
}

var CommentSet = wire.NewSet(
	wire.Struct(new(CommentService), "*"),
	wire.Bind(new(ICommentService), new(*CommentService)),
)

func CommentConvert(in db.Comment) *comment.Comment {
	return &comment.Comment{
		Id:           in.ID.Hex(),
		Text:         in.Text,
		AuthorId:     in.AuthorId,
		ReplyTo:      in.ReplyTo,
		Type:         consts.MapStringCommentType[in.Type],
		ParentId:     in.ParentId,
		FirstLevelId: in.FirstLevelId,
		UpdateAt:     in.UpdateAt.Unix(),
		CreateAt:     in.CreateAt.Unix(),
	}
}

func (s *CommentService) CreateComment(ctx context.Context, req *comment.CreateCommentReq) (resp *comment.CreateCommentResp, err error) {
	data := db.Comment{
		Text:         req.Text,
		FirstLevelId: req.FirstLevelId,
		AuthorId:     req.AuthorId,
		ReplyTo:      req.ReplyTo,
		Type:         consts.MapCommentTypeString[req.Type],
		ParentId:     req.ParentId,
	}
	if err := s.CommentModel.Insert(ctx, &data); err != nil {
		return nil, err
	}
	resp = &comment.CreateCommentResp{
		Id:      data.ID.Hex(),
		GetFish: false,
	}
	if req.GetFirstLevelId() != "" {
		return resp, nil
	}
	t, err := s.Redis.GetCtx(ctx, "commentTimes"+req.AuthorId)
	if err != nil {
		return resp, nil
	}
	r, err := s.Redis.GetCtx(ctx, "commentDate"+req.AuthorId)
	if err != nil {
		return resp, nil
	} else if r == "" {
		resp.GetFish = true
		resp.GetFishTimes = 1
		err = s.Redis.SetexCtx(ctx, "commentTimes"+req.AuthorId, "1", 86400)
		if err != nil {
			resp.GetFish = false
			return resp, nil
		}
		err = s.Redis.SetexCtx(ctx, "commentDate"+req.AuthorId, strconv.FormatInt(time.Now().Unix(), 10), 86400)
		if err != nil {
			resp.GetFish = false
			return resp, nil
		}
	} else {
		times, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return resp, nil
		}
		resp.GetFishTimes = times + 1
		date, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			return resp, nil
		}
		lastTime := time.Unix(date, 0)
		err = s.Redis.SetexCtx(ctx, "commentTimes"+req.AuthorId, strconv.FormatInt(times+1, 10), 86400)
		if err != nil {
			return resp, nil
		}
		err = s.Redis.SetexCtx(ctx, "commentDate"+req.AuthorId, strconv.FormatInt(time.Now().Unix(), 10), 86400)
		if err != nil {
			return resp, nil
		}
		if lastTime.Day() == time.Now().Day() && lastTime.Month() == time.Now().Month() && lastTime.Year() == time.Now().Year() {
			err = s.Redis.SetexCtx(ctx, "commentTimes"+req.AuthorId, strconv.FormatInt(times+1, 10), 86400)
			if err != nil {
				return resp, nil
			}
			if times >= s.Config.GetFishTimes {
				resp.GetFish = false
			} else {
				resp.GetFish = true
			}
		} else {
			err = s.Redis.SetexCtx(ctx, "commentTimes"+req.AuthorId, "1", 86400)
			if err != nil {
				return resp, nil
			}
			resp.GetFish = true
			resp.GetFishTimes = 1
		}
	}
	return resp, nil
}

func (s *CommentService) saveToHistory(ctx context.Context, historyType string, data *db.Comment) error {
	return s.HistoryModel.Insert(ctx, &db.History{
		Type: historyType,
		Data: *data,
	})
}

func (s *CommentService) UpdateComment(ctx context.Context, req *comment.UpdateCommentReq) (resp *comment.UpdateCommentResp, err error) {
	old, err := s.CommentModel.FindOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	err = s.saveToHistory(ctx, "update", old)

	if err != nil {
		return nil, err
	}
	old.Text = req.Text
	err = s.CommentModel.Update(ctx, old)
	if err != nil {
		return nil, err
	}
	return &comment.UpdateCommentResp{}, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, req *comment.DeleteCommentByIdReq) (resp *comment.DeleteCommentByIdResp, err error) {

	var data *db.Comment

	if data, err = s.CommentModel.FindOne(ctx, req.Id); err != nil {
		return nil, err
	}
	if err = s.saveToHistory(ctx, "delete", data); err != nil {

		return nil, err
	}
	if err = s.CommentModel.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	if data.FirstLevelId == "" {
		children, _, err := s.CommentModel.FindByFirstLevel(ctx, data.ID.Hex(), 0, 9999)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			var childData *db.Comment
			if childData, err = s.CommentModel.FindOne(ctx, child.ID.Hex()); err != nil {
				return nil, err
			}
			if err = s.saveToHistory(ctx, "delete", childData); err != nil {
				return nil, err
			}
			if err = s.CommentModel.Delete(ctx, child.ID.Hex()); err != nil {
				return nil, err
			}
		}
	}

	return &comment.DeleteCommentByIdResp{}, nil
}

func (s *CommentService) ListCommentByParent(ctx context.Context, req *comment.ListCommentByParentReq) (resp *comment.ListCommentByParentResp, err error) {
	var data []db.Comment
	var count int64
	if consts.MapCommentTypeString[req.Type] != "comment" {
		data, count, err = s.CommentModel.FindByParent(ctx, consts.MapCommentTypeString[req.Type], req.Id, req.OnlyFirstLevel, req.Skip, req.Limit)
	} else {
		data, count, err = s.CommentModel.FindByFirstLevel(ctx, req.Id, req.Skip, req.Limit)
	}

	if err != nil {
		return nil, err
	}

	res := comment.ListCommentByParentResp{
		Comments: make([]*comment.Comment, 0, len(data)),
		Total:    count,
	}
	for _, val := range data {
		res.Comments = append(res.Comments,
			CommentConvert(val),
		)
	}
	return &res, nil
}

func (s *CommentService) CountCommentByParent(ctx context.Context, req *comment.CountCommentByParentReq) (resp *comment.CountCommentByParentResp, err error) {
	if consts.MapCommentTypeString[req.Type] != "comment" {
		total, err := s.CommentModel.CountByParent(ctx, consts.MapCommentTypeString[req.Type], req.ParentId, req.OnlyFirstLevel)
		if err != nil {
			return nil, err
		}
		return &comment.CountCommentByParentResp{Total: total}, nil
	} else {
		total, err := s.CommentModel.CountByFirstLevel(ctx, req.ParentId)
		if err != nil {
			return nil, err
		}
		return &comment.CountCommentByParentResp{Total: total}, nil
	}
}

func (s *CommentService) RetrieveCommentById(ctx context.Context, req *comment.RetrieveCommentByIdReq) (resp *comment.RetrieveCommentByIdResp, err error) {
	ret, err := s.CommentModel.FindOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &comment.RetrieveCommentByIdResp{
		Comment: CommentConvert(*ret),
	}, nil
}

func (s *CommentService) ListCommentByAuthorIdAndType(ctx context.Context, req *comment.ListCommentByAuthorIdAndTypeReq) (resp *comment.ListCommentByAuthorIdAndTypeResp, err error) {
	data, count, err := s.CommentModel.FindByAuthorIdAndType(ctx, req.AuthorId, consts.MapCommentTypeString[req.Type], req.Skip, req.Limit)
	if err != nil {
		return nil, err
	}
	res := comment.ListCommentByAuthorIdAndTypeResp{
		Comments: make([]*comment.Comment, 0, len(data)),
		Total:    count,
	}
	for _, val := range data {
		res.Comments = append(res.Comments, CommentConvert(val))
	}
	return &res, nil
}

func (s *CommentService) ListCommentByReplyToAndType(ctx context.Context, req *comment.ListCommentByReplyToAndTypeReq) (resp *comment.ListCommentByReplyToAndTypeResp, err error) {
	data, count, err := s.CommentModel.FindByReplyToAndType(ctx, consts.MapCommentTypeString[req.Type], req.ReplyTo, req.Skip, req.Limit)
	if err != nil {
		return nil, err
	}
	res := comment.ListCommentByReplyToAndTypeResp{
		Comments: make([]*comment.Comment, 0, len(data)),
		Total:    count,
	}
	for _, val := range data {
		res.Comments = append(res.Comments,
			CommentConvert(val),
		)
	}
	return &res, nil
}
