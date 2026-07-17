package services

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	commentpb "comment/gen"
)

type CommentService struct {
	client commentpb.CommentServiceClient
	conn   *grpc.ClientConn
	token  string
}

func NewCommentService(address, sharedToken string) (*CommentService, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &CommentService{
		client: commentpb.NewCommentServiceClient(conn),
		conn:   conn,
		token:  sharedToken,
	}, nil
}

func (s *CommentService) authContext(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "x-internal-token", s.token)
}

type CampaignComment struct {
	ID         string `json:"id"`
	CampaignID string `json:"campaignId"`
	AuthorSub  string `json:"authorSub"`
	AuthorName string `json:"authorName"`
	Text       string `json:"text"`
	ParentID   string `json:"parentId,omitempty"`
	CreatedAt  int64  `json:"createdAt"`
}

func (s *CommentService) ListComments(ctx context.Context, campaignID string, offset, limit uint64) ([]CampaignComment, int64, error) {
	resp, err := s.client.ListComments(s.authContext(ctx), &commentpb.ListCommentsRequest{
		CampaignId: campaignID,
		Offset:     offset,
		Limit:      limit,
	})
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			return nil, 0, NewUnavailableError("comment service is temporarily unavailable, please try again shortly")
		}
		if status.Code(err) == codes.InvalidArgument {
			return nil, 0, NewValidationError(status.Convert(err).Message())
		}
		return nil, 0, err
	}

	items := make([]CampaignComment, len(resp.Items))
	for i, item := range resp.Items {
		items[i] = CampaignComment{
			ID:         item.Id,
			CampaignID: item.CampaignId,
			AuthorSub:  item.AuthorSub,
			AuthorName: item.AuthorName,
			Text:       item.Text,
			ParentID:   item.ParentId,
			CreatedAt:  item.CreatedAt,
		}
	}

	return items, resp.Total, nil
}

func (s *CommentService) PostComment(ctx context.Context, campaignID, authorSub, authorName, text string) (*CampaignComment, error) {
	resp, err := s.client.PostComment(s.authContext(ctx), &commentpb.PostCommentRequest{
		CampaignId: campaignID,
		AuthorSub:  authorSub,
		AuthorName: authorName,
		Text:       text,
	})
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			return nil, NewUnavailableError("comment service is temporarily unavailable, please try again shortly")
		}
		if status.Code(err) == codes.InvalidArgument {
			return nil, NewValidationError(status.Convert(err).Message())
		}
		return nil, err
	}

	return &CampaignComment{
		ID:         resp.Id,
		CampaignID: resp.CampaignId,
		AuthorSub:  resp.AuthorSub,
		AuthorName: resp.AuthorName,
		Text:       resp.Text,
		ParentID:   resp.ParentId,
		CreatedAt:  resp.CreatedAt,
	}, nil
}

func (s *CommentService) GetComment(ctx context.Context, id string) (*CampaignComment, error) {
	resp, err := s.client.GetComment(s.authContext(ctx), &commentpb.GetCommentRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, NewNotFoundError(status.Convert(err).Message())
		}
		if status.Code(err) == codes.Unavailable {
			return nil, NewUnavailableError("comment service is temporarily unavailable, please try again shortly")
		}
		return nil, err
	}

	return &CampaignComment{
		ID:         resp.Id,
		CampaignID: resp.CampaignId,
		AuthorSub:  resp.AuthorSub,
		AuthorName: resp.AuthorName,
		Text:       resp.Text,
		ParentID:   resp.ParentId,
		CreatedAt:  resp.CreatedAt,
	}, nil
}

func (s *CommentService) ReplyToComment(ctx context.Context, parentID, authorSub, authorName, text string) (*CampaignComment, error) {
	resp, err := s.client.ReplyToComment(s.authContext(ctx), &commentpb.ReplyToCommentRequest{
		ParentId:   parentID,
		AuthorSub:  authorSub,
		AuthorName: authorName,
		Text:       text,
	})
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			return nil, NewUnavailableError("comment service is temporarily unavailable, please try again shortly")
		}
		if status.Code(err) == codes.InvalidArgument {
			return nil, NewValidationError(status.Convert(err).Message())
		}
		if status.Code(err) == codes.NotFound {
			return nil, NewNotFoundError(status.Convert(err).Message())
		}
		return nil, err
	}

	return &CampaignComment{
		ID:         resp.Id,
		CampaignID: resp.CampaignId,
		AuthorSub:  resp.AuthorSub,
		AuthorName: resp.AuthorName,
		Text:       resp.Text,
		ParentID:   resp.ParentId,
		CreatedAt:  resp.CreatedAt,
	}, nil
}
