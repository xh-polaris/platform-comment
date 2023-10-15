package adaptor

import (
	"context"
	"github.com/xh-polaris/platform-comment/biz/application/service"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/platform/comment"

	"github.com/xh-polaris/platform-comment/biz/infrastructure/config"
)

type CommentServerImpl struct {
	*config.Config
	CommentService service.ICommentService
}

func (s *CommentServerImpl) CreateComment(ctx context.Context, req *comment.CreateCommentReq) (res *comment.CreateCommentResp, err error) {
	return s.CommentService.CreateComment(ctx, req)
}

func (s *CommentServerImpl) UpdateComment(ctx context.Context, req *comment.UpdateCommentReq) (res *comment.UpdateCommentResp, err error) {
	return s.CommentService.UpdateComment(ctx, req)
}

func (s *CommentServerImpl) DeleteComment(ctx context.Context, req *comment.DeleteCommentByIdReq) (res *comment.DeleteCommentByIdResp, err error) {
	return s.CommentService.DeleteComment(ctx, req)
}

func (s *CommentServerImpl) ListCommentByParent(ctx context.Context, req *comment.ListCommentByParentReq) (res *comment.ListCommentByParentResp, err error) {
	return s.CommentService.ListCommentByParent(ctx, req)

}

func (s *CommentServerImpl) CountCommentByParent(ctx context.Context, req *comment.CountCommentByParentReq) (res *comment.CountCommentByParentResp, err error) {
	return s.CommentService.CountCommentByParent(ctx, req)
}

func (s *CommentServerImpl) RetrieveCommentById(ctx context.Context, req *comment.RetrieveCommentByIdReq) (res *comment.RetrieveCommentByIdResp, err error) {
	return s.CommentService.RetrieveCommentById(ctx, req)

}

func (s *CommentServerImpl) ListCommentByAuthorIdAndType(ctx context.Context, req *comment.ListCommentByAuthorIdAndTypeReq) (res *comment.ListCommentByAuthorIdAndTypeResp, err error) {
	return s.CommentService.ListCommentByAuthorIdAndType(ctx, req)
}

func (s *CommentServerImpl) ListCommentByReplyToAndType(ctx context.Context, req *comment.ListCommentByReplyToAndTypeReq) (res *comment.ListCommentByReplyToAndTypeResp, err error) {
	return s.CommentService.ListCommentByReplyToAndType(ctx, req)
}
