package grpcserver

import (
	"context"
	"errors"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
	userphoto "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user_photo"
	userservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user"
	userphotosvc "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user_photo"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userHandler struct {
	pb.UnimplementedUserServiceServer
	userService      userservice.UserService
	userPhotoService userphotosvc.UserPhotoService
}

func newUserHandler(userService userservice.UserService, userPhotoService userphotosvc.UserPhotoService) *userHandler {
	return &userHandler{
		userService:      userService,
		userPhotoService: userPhotoService,
	}
}

// photoSvcErrToGRPC maps user_photo service errors to gRPC status errors.
// The service knows nothing about gRPC — this mapping lives at the transport boundary.
func photoSvcErrToGRPC(err error) error {
	switch {
	case errors.Is(err, userphotosvc.ErrPhotoNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, userphotosvc.ErrPhotoNotOwned):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, userphotosvc.ErrFileNotFound):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, userphotosvc.ErrFileNotOwned):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, userphotosvc.ErrInvalidFileBucket):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, userphotosvc.ErrFileNotReady):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
}

// photoToProto converts a domain UserPhoto to its protobuf representation.
// Lives here because it belongs to the transport boundary, not the service.
func photoToProto(p *userphoto.UserPhoto) *pb.UserPhoto {
	if p == nil {
		return nil
	}
	return &pb.UserPhoto{
		Id:        p.GetID(),
		FileId:    p.GetFileID(),
		CreatedAt: p.GetCreatedAt().Format(time.RFC3339),
	}
}

func photosToProto(photos []*userphoto.UserPhoto) []*pb.UserPhoto {
	result := make([]*pb.UserPhoto, 0, len(photos))
	for _, p := range photos {
		result = append(result, photoToProto(p))
	}
	return result
}

func (h *userHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	utils.SanitizeCreateUserRequest(req)
	return h.userService.CreateUser(ctx, req)
}

func (h *userHandler) ConfirmUser(ctx context.Context, req *pb.ConfirmUserRequest) (*pb.ConfirmUserResponse, error) {
	utils.SanitizeConfirmUserRequest(req)
	return h.userService.ConfirmUser(ctx, req)
}

func (h *userHandler) UpdateUserEmail(ctx context.Context, req *pb.UpdateUserEmailRequest) (*pb.UpdateUserEmailResponse, error) {
	utils.SanitizeUpdateUserEmailRequest(req)
	return h.userService.UpdateUserEmail(ctx, req)
}

func (h *userHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	utils.SanitizeGetUserByIDRequest(req)
	resp, err := h.userService.GetUserByID(ctx, req)
	if err != nil || resp == nil || resp.User == nil || !req.IncludeAvatar {
		return resp, err
	}
	photo, _ := h.userPhotoService.GetLatestPhoto(ctx, resp.User.Id)
	resp.User.Avatar = photoToProto(photo)
	return resp, nil
}

func (h *userHandler) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	utils.SanitizeGetUserByUsernameRequest(req)
	resp, err := h.userService.GetUserByUsername(ctx, req)
	if err != nil || resp == nil || resp.User == nil || !req.IncludeAvatar {
		return resp, err
	}
	photo, _ := h.userPhotoService.GetLatestPhoto(ctx, resp.User.Id)
	resp.User.Avatar = photoToProto(photo)
	return resp, nil
}

func (h *userHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	utils.SanitizeGetUserByEmailRequest(req)
	resp, err := h.userService.GetUserByEmail(ctx, req)
	if err != nil || resp == nil || resp.User == nil || !req.IncludeAvatar {
		return resp, err
	}
	photo, _ := h.userPhotoService.GetLatestPhoto(ctx, resp.User.Id)
	resp.User.Avatar = photoToProto(photo)
	return resp, nil
}

func (h *userHandler) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	utils.SanitizeGetUsersRequest(req)
	resp, err := h.userService.GetUsers(ctx, req)
	if err != nil || resp == nil || !req.IncludeAvatar || len(resp.Users) == 0 {
		return resp, err
	}

	ids := make([]string, len(resp.Users))
	for i, u := range resp.Users {
		ids[i] = u.Id
	}
	avatars, _ := h.userPhotoService.GetLatestPhotosByUserIDs(ctx, ids)
	if avatars != nil {
		for _, u := range resp.Users {
			u.Avatar = photoToProto(avatars[u.Id])
		}
	}
	return resp, nil
}

// GetMeByUser always returns the avatar.
func (h *userHandler) GetMeByUser(ctx context.Context, req *pb.GetMeByUserRequest) (*pb.GetMeByUserResponse, error) {
	resp, err := h.userService.GetMeByUser(ctx, req)
	if err != nil || resp == nil || resp.User == nil {
		return resp, err
	}
	photo, _ := h.userPhotoService.GetLatestPhoto(ctx, resp.User.Id)
	resp.User.Avatar = photoToProto(photo)
	return resp, nil
}

func (h *userHandler) UpdateMeByUser(ctx context.Context, req *pb.UpdateMeByUserRequest) (*pb.UpdateMeByUserResponse, error) {
	utils.SanitizeUpdateMeByUserRequest(req)
	return h.userService.UpdateMeByUser(ctx, req)
}

func (h *userHandler) UpdateMyPrivacySettingsByUserRequest(ctx context.Context, req *pb.UpdateMyPrivacySettingsByUserRequest) (*pb.UpdateMyPrivacySettingsByUserResponse, error) {
	utils.SanitizeUpdateMyPrivacySettingsByUserRequest(req)
	return h.userService.UpdateMyPrivacySettingsByUserRequest(ctx, req)
}

func (h *userHandler) DeleteMeByUser(ctx context.Context, req *pb.DeleteMeByUserRequest) (*pb.DeleteMeByUserResponse, error) {
	return h.userService.DeleteMeByUser(ctx, req)
}

func (h *userHandler) AddUserPhoto(ctx context.Context, req *pb.AddUserPhotoRequest) (*pb.AddUserPhotoResponse, error) {
	utils.SanitizeAddUserPhotoRequest(req)

	claims, _ := ctxmetadata.GetUserClaimsFromContext(ctx)
	if claims == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	photo, err := h.userPhotoService.AddPhoto(ctx, claims.Id, req.FileId)
	if err != nil {
		return nil, photoSvcErrToGRPC(err)
	}
	return &pb.AddUserPhotoResponse{Photo: photoToProto(photo)}, nil
}

func (h *userHandler) GetUserPhotos(ctx context.Context, req *pb.GetUserPhotosRequest) (*pb.GetUserPhotosResponse, error) {
	utils.SanitizeGetUserPhotosRequest(req)

	photos, total, err := h.userPhotoService.GetPhotos(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, photoSvcErrToGRPC(err)
	}
	return &pb.GetUserPhotosResponse{Photos: photosToProto(photos), Total: total}, nil
}

func (h *userHandler) DeleteUserPhoto(ctx context.Context, req *pb.DeleteUserPhotoRequest) (*pb.DeleteUserPhotoResponse, error) {
	utils.SanitizeDeleteUserPhotoRequest(req)

	claims, _ := ctxmetadata.GetUserClaimsFromContext(ctx)
	if claims == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	if err := h.userPhotoService.DeletePhoto(ctx, req.PhotoId, claims.Id); err != nil {
		return nil, photoSvcErrToGRPC(err)
	}
	return &pb.DeleteUserPhotoResponse{}, nil
}

func (h *userHandler) BanUser(ctx context.Context, req *pb.BanUserRequest) (*pb.BanUserResponse, error) {
	utils.SanitizeBanUserRequest(req)
	return nil, status.Error(codes.Unimplemented, "BanUser not implemented")
}

func (h *userHandler) UnbanUser(ctx context.Context, req *pb.UnbanUserRequest) (*pb.UnbanUserResponse, error) {
	utils.SanitizeUnbanUserRequest(req)
	return nil, status.Error(codes.Unimplemented, "UnbanUser not implemented")
}
