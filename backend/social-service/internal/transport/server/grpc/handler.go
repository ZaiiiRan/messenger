package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/social/v1"
	socialservice "github.com/ZaiiiRan/messenger/backend/social-service/internal/services/social_service"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/utils"
)

type socialHandler struct {
	pb.UnimplementedSocialServiceServer
	socialService socialservice.SocialService
}

func newSocialHandler(socialService socialservice.SocialService) *socialHandler {
	return &socialHandler{
		socialService: socialService,
	}
}

func (h *socialHandler) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	utils.SanitizeGetUserByIdRequest(req)
	return h.socialService.GetUserByID(ctx, req)
}

func (h *socialHandler) GetUsersByIds(ctx context.Context, req *pb.GetUsersByIdsRequest) (*pb.GetUsersByIdsResponse, error) {
	utils.SanitizeGetUsersByIdsRequest(req)
	return h.socialService.GetUsersByIDs(ctx, req)
}

func (h *socialHandler) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	utils.SanitizeGetUserByUsernameRequest(req)
	return h.socialService.GetUserByUsername(ctx, req)
}

func (h *socialHandler) GetUsersByUsernames(ctx context.Context, req *pb.GetUsersByUsernamesRequest) (*pb.GetUsersByUsernamesResponse, error) {
	utils.SanitizeGetUsersByUsernamesRequest(req)
	return h.socialService.GetUsersByUsernames(ctx, req)
}

func (h *socialHandler) AddUsersToFriends(ctx context.Context, req *pb.AddUsersToFriendsRequest) (*pb.AddUsersToFriendsResponse, error) {
	utils.SanitizeAddUsersToFriendsRequest(req)
	return h.socialService.AddUsersToFriends(ctx, req)
}

func (h *socialHandler) RemoveUsersFromFriends(ctx context.Context, req *pb.RemoveUsersFromFriendsRequest) (*pb.RemoveUsersFromFriendsResponse, error) {
	utils.SanitizeRemoveUsersFromFriendsRequest(req)
	return h.socialService.RemoveUsersFromFriends(ctx, req)
}

func (h *socialHandler) BlockUsers(ctx context.Context, req *pb.BlockUsersRequest) (*pb.BlockUsersResponse, error) {
	utils.SanitizeBlockUsersRequest(req)
	return h.socialService.BlockUsers(ctx, req)
}

func (h *socialHandler) UnblockUsers(ctx context.Context, req *pb.UnblockUsersRequest) (*pb.UnblockUsersResponse, error) {
	utils.SanitizeUnblockUsersRequest(req)
	return h.socialService.UnblockUsers(ctx, req)
}
