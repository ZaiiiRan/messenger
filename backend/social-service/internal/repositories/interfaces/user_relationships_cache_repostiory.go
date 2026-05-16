package interfaces

import (
	"context"

	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
)

type UserRelationshipsCacheRepository interface {
	SetUserRelationship(ctx context.Context, ur *userrelationship.UserRelationship) error
	GetUserRelationship(ctx context.Context, firstUserId, secondUserId string) (*userrelationship.UserRelationship, error)
	DelUserRelationship(ctx context.Context, firstUserId, secondUserId string) error

	SetUserRelationshipsList(ctx context.Context, query *models.QueryUserRelationshipsDal, urs []*userrelationship.UserRelationship) error
	GetUserRelationshipsList(ctx context.Context, query *models.QueryUserRelationshipsDal) ([]*userrelationship.UserRelationship, error)
	InvalidateUserRelationshipsLists(ctx context.Context, userId string) error
}
