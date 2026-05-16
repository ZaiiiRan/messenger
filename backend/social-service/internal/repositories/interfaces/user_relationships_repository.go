package interfaces

import (
	"context"

	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
)

type UserRelationshipsRepository interface {
	CreateUserRelationships(ctx context.Context, urs []*userrelationship.UserRelationship) error
	UpdateUserRelationships(ctx context.Context, urs []*userrelationship.UserRelationship) error
	DeleteUserRelationships(ctx context.Context, urs []*userrelationship.UserRelationship) error
	QueryUserRelationships(ctx context.Context, query *models.QueryUserRelationshipsDal, forUpdate bool) ([]*userrelationship.UserRelationship, error)
}
