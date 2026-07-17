package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commentpb "comment/gen"
)

type CommentServer struct {
	commentpb.UnimplementedCommentServiceServer
}

func NewCommentServer() *CommentServer {
	return &CommentServer{}
}

func (s *CommentServer) PostComment(ctx context.Context, req *commentpb.PostCommentRequest) (*commentpb.Comment, error) {
	return nil, status.Errorf(codes.Unimplemented, "PostComment not implemented yet")
}

func (s *CommentServer) ListComments(ctx context.Context, req *commentpb.ListCommentsRequest) (*commentpb.ListCommentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ListComments not implemented yet")
}

func (s *CommentServer) ReplyToComment(ctx context.Context, req *commentpb.ReplyToCommentRequest) (*commentpb.Comment, error) {
	return nil, status.Errorf(codes.Unimplemented, "ReplyToComment not implemented yet")
}
