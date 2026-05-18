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

func (h *socialHandler) GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (*pb.GetFriendsResponse, error) {
	utils.SanitizeGetFriendsRequest(req)
	return h.socialService.GetFriends(ctx, req)
}

func (h *socialHandler) GetIncomingFriendRequests(ctx context.Context, req *pb.GetIncomingFriendRequestsRequest) (*pb.GetIncomingFriendRequestsResponse, error) {
	utils.SanitizeGetIncomingFriendRequestsRequest(req)
	return h.socialService.GetIncomingFriendRequests(ctx, req)
}

func (h *socialHandler) GetOutgoingFriendRequests(ctx context.Context, req *pb.GetOutgoingFriendRequestsRequest) (*pb.GetOutgoingFriendRequestsResponse, error) {
	utils.SanitizeGetOutgoingFriendRequestsRequest(req)
	return h.socialService.GetOutgoingFriendRequests(ctx, req)
}

func (h *socialHandler) GetBlockedUsers(ctx context.Context, req *pb.GetBlockedUsersRequest) (*pb.GetBlockedUsersResponse, error) {
	utils.SanitizeGetBlockedUsersRequest(req)
	return h.socialService.GetBlockedUsers(ctx, req)
}

func (h *socialHandler) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	utils.SanitizeSearchUsersRequest(req)
	return h.socialService.SearchUsers(ctx, req)
}
