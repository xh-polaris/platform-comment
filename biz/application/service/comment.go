package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/config"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/data/db"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/mapper"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/platform/comment"
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
}

var CommentSet = wire.NewSet(
	wire.Struct(new(CommentService), "*"),
	wire.Bind(new(ICommentService), new(*CommentService)),
)

func CommentConvert(in db.Comment) *comment.Comment {
	return &comment.Comment{
		Id:       in.ID.Hex(),
		Text:     in.Text,
		AuthorId: in.AuthorId,
		ReplyTo:  in.ReplyTo,
		Type:     in.Type,
		ParentId: in.ParentId,
		UpdateAt: in.UpdateAt.Unix(),
		CreateAt: in.CreateAt.Unix(),
	}
}

func (s *CommentService) CreateComment(ctx context.Context, req *comment.CreateCommentReq) (resp *comment.CreateCommentResp, err error) {
	data := db.Comment{
		Text:     req.Text,
		AuthorId: req.AuthorId,
		ReplyTo:  req.ReplyTo,
		Type:     req.Type,
		ParentId: req.ParentId,
	}
	if err := s.CommentModel.Insert(ctx, &data); err != nil {
		return nil, err
	}
	return &comment.CreateCommentResp{
		Id: data.ID.Hex(),
	}, nil
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

	return &comment.DeleteCommentByIdResp{}, nil
}

func (s *CommentService) ListCommentByParent(ctx context.Context, req *comment.ListCommentByParentReq) (resp *comment.ListCommentByParentResp, err error) {
	data, count, err := s.CommentModel.FindByParent(ctx, req.Type, req.ParentId, req.Skip, req.Limit)
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
	total, err := s.CommentModel.CountByParent(ctx, req.Type, req.ParentId)
	if err != nil {
		return nil, err
	}
	return &comment.CountCommentByParentResp{Total: total}, nil
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
	data, count, err := s.CommentModel.FindByAuthorIdAndType(ctx, req.AuthorId, req.Type, req.Skip, req.Limit)
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
	data, count, err := s.CommentModel.FindByReplyToAndType(ctx, req.Type, req.ReplyTo, req.Skip, req.Limit)
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