package fileservice

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ZaiiiRan/messenger/backend/file-service/gen/go/file/v1"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
	filevariant "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file_variant"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

type FileService interface {
	CreateFilesByUser(ctx context.Context, req *pb.CreateFilesByUserRequest) (*pb.CreateFilesByUserResponse, error)
	ConfirmUploadsByUser(ctx context.Context, req *pb.ConfirmUploadsByUserRequest) (*pb.ConfirmUploadsByUserResponse, error)
	GetFilesByUser(ctx context.Context, req *pb.GetFilesByUserRequest) (*pb.GetFilesByUserResponse, error)
	GetFiles(ctx context.Context, req *pb.GetFilesRequest) (*pb.GetFilesResponse, error)
	DeleteFiles(ctx context.Context, req *pb.DeleteFilesRequest) (*pb.DeleteFilesResponse, error)
	UpdateFiles(ctx context.Context, req *pb.UpdateFilesRequest) (*pb.UpdateFilesResponse, error)
}

type service struct {
	fileDataProvider        *fileDataProvider
	fileVariantDataProvider *fileVariantDataProvider
	fileStorage             *fileStorage
	log                     *zap.SugaredLogger
}

func (s *service) CreateFilesByUser(ctx context.Context, req *pb.CreateFilesByUserRequest) (*pb.CreateFilesByUserResponse, error) {
	l := s.log.With("op", "create_files_by_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	actor, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.Files) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrFileIdsIsRequired.Error())
	}
	if len(req.Files) > 10 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyFiles.Error())
	}

	type reqItem struct {
		purpose  file.FilePurpose
		mimeType string
		size     int64
	}
	items := make([]reqItem, len(req.Files))
	for i, fr := range req.Files {
		purpose, err := file.FilePurposeFromString(fr.FilePurpose)
		if err != nil {
			return nil, grpcstatus.Errorf(codes.InvalidArgument, "files[%d]: invalid file_purpose %q", i, fr.FilePurpose)
		}
		if err := filevariant.ValidateMimeType(fr.MimeType, purpose); err != nil {
			return nil, grpcstatus.Errorf(codes.InvalidArgument, "files[%d]: %s", i, err.Error())
		}
		if fr.Size <= 0 {
			return nil, grpcstatus.Errorf(codes.InvalidArgument, "files[%d]: size must be positive", i)
		}
		items[i] = reqItem{purpose: purpose, mimeType: fr.MimeType, size: fr.Size}
	}

	domainFiles := make([]*file.File, len(items))
	for i := range items {
		domainFiles[i] = file.New(actor.Id)
	}

	uow := s.fileDataProvider.newUOW()
	defer uow.Close()
	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("file.create_files_by_user_failed.begin_transaction_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	if err := s.fileDataProvider.createFiles(ctx, domainFiles, uow); err != nil {
		l.Errorw("file.create_files_by_user_failed.create_files_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	variants := make([]*filevariant.FileVariant, len(domainFiles))
	type uploadInfo struct {
		url      string
		formData map[string]string
	}
	uploads := make([]uploadInfo, len(domainFiles))

	for i, f := range domainFiles {
		item := items[i]
		bucket, presign := s.purposeStorage(item.purpose)
		s3Key := fmt.Sprintf("%s/%s/original", item.purpose.String(), f.GetID())

		fv, _ := filevariant.New(f.GetID(), bucket, s3Key, item.mimeType)
		variants[i] = fv

		uploadURL, formData, err := presign(ctx, s3Key, item.mimeType, item.size)
		if err != nil {
			l.Errorw("file.create_files_by_user_failed.presign_upload_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
		uploads[i] = uploadInfo{url: uploadURL, formData: formData}
	}

	if err := s.fileVariantDataProvider.createFileVariants(ctx, variants, uow); err != nil {
		l.Errorw("file.create_files_by_user_failed.create_file_variants_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("file.create_files_by_user_failed.commit_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	results := make([]*pb.CreateFileResult, len(domainFiles))
	for i, f := range domainFiles {
		pbVariant, err := fileVariantToProto(variants[i], nil, nil, false)
		if err != nil {
			l.Errorw("file.create_files_by_user_failed.map_variant_to_proto_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
		results[i] = &pb.CreateFileResult{
			File:         fileToProto(f, []*pb.FileVariant{pbVariant}, false),
			UrlForUpload: uploads[i].url,
			FormData:     uploads[i].formData,
		}
	}

	l.Infow("file.create_files_by_user.success", "count", len(results), "owner_id", actor.Id)
	return &pb.CreateFilesByUserResponse{Files: results}, nil
}

func (s *service) ConfirmUploadsByUser(ctx context.Context, req *pb.ConfirmUploadsByUserRequest) (*pb.ConfirmUploadsByUserResponse, error) {
	l := s.log.With("op", "confirm_uploads_by_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	actor, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.FileIds) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrFileIdsIsRequired.Error())
	}
	if len(req.FileIds) > 10 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyFiles.Error())
	}

	uow := s.fileDataProvider.newUOW()
	defer uow.Close()

	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("file.confirm_uploads_by_user_failed.begin_transaction_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	var files []*file.File
	if len(req.FileIds) == 1 {
		f, err := s.fileDataProvider.getFileByIDLocked(ctx, req.FileIds[0], uow)
		if err != nil {
			l.Errorw("file.confirm_uploads_by_user_failed.get_file_locked_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
		if f != nil && f.GetOwnerID() == actor.Id {
			files = []*file.File{f}
		}
	} else {
		query := models.NewQueryFilesDal(
			models.FileFilterDal{
				Ids:      req.FileIds,
				OwnerIds: []string{actor.Id},
			},
			1, len(req.FileIds),
		)
		files, err = s.fileDataProvider.getFilesLocked(ctx, query, false, uow)
		if err != nil {
			l.Errorw("file.confirm_uploads_by_user_failed.get_files_locked_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
	}

	now := time.Now()
	toUpdate := make([]*file.File, 0, len(files))
	for _, f := range files {
		if err := f.SetStatus(file.FileStatusUploaded); err != nil {
			continue
		}
		f.SetUpdatedAt(&now)
		toUpdate = append(toUpdate, f)
	}

	if len(toUpdate) > 0 {
		if err := s.fileDataProvider.updateFiles(ctx, toUpdate, uow); err != nil {
			l.Errorw("file.confirm_uploads_by_user_failed.update_files_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("file.confirm_uploads_by_user_failed.commit_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	readUow := s.fileVariantDataProvider.newUOW()
	defer readUow.Close()

	pbFiles := make([]*pb.File, 0, len(toUpdate))
	for _, f := range toUpdate {
		fvs, err := s.fileVariantDataProvider.getFileVariantsByFileID(ctx, f.GetID(), readUow)
		if err != nil {
			l.Errorw("file.confirm_uploads_by_user_failed.get_variants_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		pbVariants := make([]*pb.FileVariant, 0, len(fvs))
		for _, fv := range fvs {
			downloadURL, err := s.bucketPresignDownload(ctx, fv.GetBucket(), fv.GetS3Key())
			if err != nil {
				l.Errorw("file.confirm_uploads_by_user_failed.presign_download_error", "err", err)
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			pbVariant, err := fileVariantToProto(fv, &downloadURL, nil, false)
			if err != nil {
				l.Errorw("file.confirm_uploads_by_user_failed.map_variant_error", "err", err)
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			pbVariants = append(pbVariants, pbVariant)
		}

		pbFiles = append(pbFiles, fileToProto(f, pbVariants, false))
	}

	l.Infow("file.confirm_uploads_by_user.success", "confirmed", len(pbFiles), "owner_id", actor.Id)
	return &pb.ConfirmUploadsByUserResponse{File: pbFiles}, nil
}

func (s *service) GetFilesByUser(ctx context.Context, req *pb.GetFilesByUserRequest) (*pb.GetFilesByUserResponse, error) {
	l := s.log.With("op", "get_files_by_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	actor, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.FileIds) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrFileIdsIsRequired.Error())
	}
	if len(req.FileIds) > 10 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyFiles.Error())
	}

	uow := s.fileDataProvider.newUOW()
	defer uow.Close()

	var files []*file.File
	if len(req.FileIds) == 1 {
		f, err := s.fileDataProvider.getFileByID(ctx, req.FileIds[0], uow)
		if err != nil {
			l.Errorw("file.get_files_by_user_failed.get_file_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
		if f != nil && f.GetOwnerID() == actor.Id {
			files = []*file.File{f}
		}
	} else {
		query := models.NewQueryFilesDal(
			models.FileFilterDal{
				Ids:      req.FileIds,
				OwnerIds: []string{actor.Id},
			},
			1, len(req.FileIds),
		)
		files, err = s.fileDataProvider.getFiles(ctx, query, uow)
		if err != nil {
			l.Errorw("file.get_files_by_user_failed.get_files_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
	}

	pbFiles := make([]*pb.File, 0, len(files))
	for _, f := range files {
		fvs, err := s.fileVariantDataProvider.getFileVariantsByFileID(ctx, f.GetID(), uow)
		if err != nil {
			l.Errorw("file.get_files_by_user_failed.get_variants_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		pbVariants := make([]*pb.FileVariant, 0, len(fvs))
		for _, fv := range fvs {
			downloadURL, err := s.bucketPresignDownload(ctx, fv.GetBucket(), fv.GetS3Key())
			if err != nil {
				l.Errorw("file.get_files_by_user_failed.presign_download_error", "err", err)
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			pbVariant, err := fileVariantToProto(fv, &downloadURL, nil, false)
			if err != nil {
				l.Errorw("file.get_files_by_user_failed.map_variant_error", "err", err)
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			pbVariants = append(pbVariants, pbVariant)
		}

		pbFiles = append(pbFiles, fileToProto(f, pbVariants, false))
	}

	l.Infow("file.get_files_by_user.success", "count", len(pbFiles), "owner_id", actor.Id)
	return &pb.GetFilesByUserResponse{File: pbFiles}, nil
}

func (s *service) purposeStorage(
	purpose file.FilePurpose,
) (filevariant.BucketType, func(ctx context.Context, key, contentType string, maxSize int64) (string, map[string]string, error)) {
	if purpose == file.PurposeUserPhoto {
		return filevariant.UserPhotosBucket, s.fileStorage.userPhotoPresignedUpload
	}
	return filevariant.ChatFilesBucket, s.fileStorage.chatFilePresignedUpload
}

func (s *service) bucketPresignDownload(ctx context.Context, bucket filevariant.BucketType, key string) (string, error) {
	if bucket == filevariant.UserPhotosBucket {
		return s.fileStorage.userPhotoPresignedDownload(ctx, key)
	}
	return s.fileStorage.chatFilePresignedDownload(ctx, key)
}

func (s *service) getAndCheckActor(ctx context.Context) (*jwt.UserClaims, error) {
	claims, _ := ctxmetadata.GetUserClaimsFromContext(ctx)
	if claims == nil || claims.IsDeleted {
		return nil, grpcstatus.Error(codes.Unauthenticated, commonerror.ErrUnauthorized.Error())
	}

	if !claims.IsConfirmed || claims.IsPermanentlyBanned || claims.IsTemporarilyBanned {
		return nil, grpcstatus.Error(codes.PermissionDenied, commonerror.ErrPermissionDenied.Error())
	}
	return claims, nil
}
