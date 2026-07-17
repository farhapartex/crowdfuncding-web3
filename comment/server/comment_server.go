package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	commentpb "comment/gen"
	"comment/models"
)

type CommentServer struct {
	commentpb.UnimplementedCommentServiceServer
	db *gorm.DB
}

func NewCommentServer(db *gorm.DB) *CommentServer {
	return &CommentServer{db: db}
}

func (s *CommentServer) PostComment(ctx context.Context, req *commentpb.PostCommentRequest) (*commentpb.Comment, error) {
	return nil, status.Errorf(codes.Unimplemented, "PostComment not implemented yet")
}

func (s *CommentServer) ListComments(ctx context.Context, req *commentpb.ListCommentsRequest) (*commentpb.ListCommentsResponse, error) {
	if req.CampaignId == "" {
		return nil, status.Error(codes.InvalidArgument, "campaign_id is required")
	}

	limit := req.Limit
	if limit == 0 || limit > 100 {
		limit = 20
	}

	comments, total, err := models.ListCommentsByCampaign(s.db, req.CampaignId, req.Offset, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list comments: %v", err)
	}

	items := make([]*commentpb.Comment, len(comments))
	for i, c := range comments {
		parentID := ""
		if c.ParentID != nil {
			parentID = *c.ParentID
		}

		items[i] = &commentpb.Comment{
			Id:         c.ID,
			CampaignId: c.CampaignID,
			AuthorSub:  c.AuthorSub,
			AuthorName: c.AuthorName,
			Text:       c.Text,
			ParentId:   parentID,
			CreatedAt:  c.CreatedAt.Unix(),
		}
	}

	return &commentpb.ListCommentsResponse{
		Items: items,
		Total: total,
	}, nil
}

func (s *CommentServer) ReplyToComment(ctx context.Context, req *commentpb.ReplyToCommentRequest) (*commentpb.Comment, error) {
	return nil, status.Errorf(codes.Unimplemented, "ReplyToComment not implemented yet")
}
