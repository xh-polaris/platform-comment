package consts

import "github.com/xh-polaris/service-idl-gen-go/kitex_gen/platform/comment"

var (
	MapCommentTypeString = map[comment.CommentType]string{
		comment.CommentType_CommentType_Unknown: "unknown",
		comment.CommentType_CommentType_Comment: "comment",
		comment.CommentType_CommentType_Post:    "post",
		comment.CommentType_CommentType_Moment:  "moment",
	}
	MapStringCommentType = map[string]comment.CommentType{
		"unknown": comment.CommentType_CommentType_Unknown,
		"comment": comment.CommentType_CommentType_Comment,
		"post":    comment.CommentType_CommentType_Post,
		"moment":  comment.CommentType_CommentType_Moment,
	}
)
