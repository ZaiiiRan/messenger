package interfaces

import (
	"context"

	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
)

type UserRelationshipsRepository interface {
	CreateUserRelationship(ctx context.Context, ur *userrelationship.UserRelationship) error
	UpdateUserRelationship(ctx context.Context, ur *userrelationship.UserRelationship) error
	DeleteUserRelationship(ctx context.Context, ur *userrelationship.UserRelationship) error
	QueryUserRelationships(ctx context.Context, query *models.QueryUserRelationshipsDal, forUpdate bool) ([]*userrelationship.UserRelationship, error)
}
