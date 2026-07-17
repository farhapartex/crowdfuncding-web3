package server

import (
	"context"
	"strings"

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
	if req.CampaignId == "" {
		return nil, status.Error(codes.InvalidArgument, "campaign_id is required")
	}
	if req.AuthorSub == "" {
		return nil, status.Error(codes.InvalidArgument, "author_sub is required")
	}
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return nil, status.Error(codes.InvalidArgument, "text is required")
	}

	comment, err := models.CreateComment(s.db, req.CampaignId, req.AuthorSub, req.AuthorName, text, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create comment: %v", err)
	}

	return &commentpb.Comment{
		Id:         comment.ID,
		CampaignId: comment.CampaignID,
		AuthorSub:  comment.AuthorSub,
		AuthorName: comment.AuthorName,
		Text:       comment.Text,
		CreatedAt:  comment.CreatedAt.Unix(),
	}, nil
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
	if req.ParentId == "" {
		return nil, status.Error(codes.InvalidArgument, "parent_id is required")
	}
	if req.AuthorSub == "" {
		return nil, status.Error(codes.InvalidArgument, "author_sub is required")
	}
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return nil, status.Error(codes.InvalidArgument, "text is required")
	}

	parent, err := models.GetCommentByID(s.db, req.ParentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to look up parent comment: %v", err)
	}
	if parent == nil {
		return nil, status.Error(codes.NotFound, "parent comment not found")
	}

	comment, err := models.CreateComment(s.db, parent.CampaignID, req.AuthorSub, req.AuthorName, text, parent.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create reply: %v", err)
	}

	return &commentpb.Comment{
		Id:         comment.ID,
		CampaignId: comment.CampaignID,
		AuthorSub:  comment.AuthorSub,
		AuthorName: comment.AuthorName,
		Text:       comment.Text,
		ParentId:   *comment.ParentID,
		CreatedAt:  comment.CreatedAt.Unix(),
	}, nil
}

func (s *CommentServer) GetComment(ctx context.Context, req *commentpb.GetCommentRequest) (*commentpb.Comment, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	comment, err := models.GetCommentByID(s.db, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to look up comment: %v", err)
	}
	if comment == nil {
		return nil, status.Error(codes.NotFound, "comment not found")
	}

	parentID := ""
	if comment.ParentID != nil {
		parentID = *comment.ParentID
	}

	return &commentpb.Comment{
		Id:         comment.ID,
		CampaignId: comment.CampaignID,
		AuthorSub:  comment.AuthorSub,
		AuthorName: comment.AuthorName,
		Text:       comment.Text,
		ParentId:   parentID,
		CreatedAt:  comment.CreatedAt.Unix(),
	}, nil
}
