//go:build wireinject
// +build wireinject

package provider

import (
	"github.com/google/wire"
	"github.com/xh-polaris/platform-comment/biz/adaptor"
)

func NewCommentServerImpl() (*adaptor.CommentServerImpl, error) {
	wire.Build(
		wire.Struct(new(adaptor.CommentServerImpl), "*"),
		AllProvider,
	)
	return nil, nil
}
