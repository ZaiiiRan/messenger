package fileservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/ZaiiiRan/messenger/backend/file-service/gen/go/file/v1"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
	filevariant "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file_variant"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/unitofwork/postgres"
	filedeletionfroms3tasksservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_deletion_from_s3_tasks_service"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/redis"
	s3transport "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/s3"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/validationerror"
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
	DeleteNotInUseFiles(ctx context.Context, workerID string, batchSize int, uow *uow.UnitOfWork) error
	DeleteFilesByIDs(ctx context.Context, workerID string, fileIDs []string, uow *uow.UnitOfWork) error
}

type service struct {
	fileDataProvider               *fileDataProvider
	fileVariantDataProvider        *fileVariantDataProvider
	fileStorage                    *fileStorage
	fileDeletionFromS3TasksService filedeletionfroms3tasksservice.FileDeletionFromS3TasksService
	log                            *zap.SugaredLogger
}

func New(
	pg *postgres.PostgresClient,
	redisClient *redis.RedisClient,
	userPhotosS3 *s3transport.S3Client,
	chatFilesS3 *s3transport.S3Client,
	fileDeletionTasksSvc filedeletionfroms3tasksservice.FileDeletionFromS3TasksService,
	log *zap.SugaredLogger,
) FileService {
	return &service{
		fileDataProvider:               newFileDataProvider(pg, redisClient),
		fileVariantDataProvider:        newFileVariantDataProvider(pg, redisClient),
		fileStorage:                    newFileStorage(userPhotosS3, chatFilesS3),
		fileDeletionFromS3TasksService: fileDeletionTasksSvc,
		log:                            log,
	}
}

type variantAddParsed struct {
	variantType filevariant.VariantType
	bucket      filevariant.BucketType
	s3Key       string
	meta        filevariant.Meta
	mimeType    string
	size        int64
}

type variantUpdateParsed struct {
	id          string
	variantType *filevariant.VariantType
	meta        *filevariant.Meta
}

type fileParsed struct {
	id               string
	purpose          *file.FilePurpose
	newStatus        *file.FileStatus
	deleteVariantIDs map[string]struct{}
	updateVariants   []variantUpdateParsed
	addVariants      []variantAddParsed
}

type updateMutations struct {
	filesToUpdate    []*file.File
	filesToDelete    []*file.File
	variantsToDelete []*filevariant.FileVariant
	variantsToUpdate []*filevariant.FileVariant
	variantsToCreate []*filevariant.FileVariant
	s3KeysByBucket   map[filevariant.BucketType][]string
}

type fileResult struct {
	f             *file.File
	finalVariants []*filevariant.FileVariant
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
		if f == nil || f.GetOwnerID() != actor.Id {
			return nil, grpcstatus.Error(codes.NotFound, ErrFileNotFound.Error())
		}
		files = []*file.File{f}
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
	return &pb.ConfirmUploadsByUserResponse{Files: pbFiles}, nil
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
		if f == nil || f.GetOwnerID() != actor.Id {
			return nil, grpcstatus.Error(codes.NotFound, ErrFileNotFound.Error())
		}
		files = []*file.File{f}
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
	return &pb.GetFilesByUserResponse{Files: pbFiles}, nil
}

func (s *service) GetFiles(ctx context.Context, req *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	l := s.log.With("op", "get_files", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(req.FileIds) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrFileIdsIsRequired.Error())
	}
	if len(req.FileIds) > 1000 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyFiles.Error())
	}

	uow := s.fileDataProvider.newUOW()
	defer uow.Close()

	var files []*file.File
	if len(req.FileIds) == 1 {
		f, err := s.fileDataProvider.getFileByID(ctx, req.FileIds[0], uow)
		if err != nil {
			l.Errorw("file.get_files_failed.get_file_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
		if f == nil {
			return nil, grpcstatus.Error(codes.NotFound, ErrFileNotFound.Error())
		}
		files = []*file.File{f}
	} else {
		query := models.NewQueryFilesDal(
			models.FileFilterDal{Ids: req.FileIds},
			1, len(req.FileIds),
		)
		f, err := s.fileDataProvider.getFiles(ctx, query, uow)
		if err != nil {
			l.Errorw("file.get_files_failed.get_files_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
		if f != nil {
			files = f
		}
	}

	pbFiles := make([]*pb.File, 0, len(files))
	for _, f := range files {
		fvs, err := s.fileVariantDataProvider.getFileVariantsByFileID(ctx, f.GetID(), uow)
		if err != nil {
			l.Errorw("file.get_files_failed.get_variants_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		pbVariants := make([]*pb.FileVariant, 0, len(fvs))
		for _, fv := range fvs {
			downloadURL, err := s.bucketPresignDownload(ctx, fv.GetBucket(), fv.GetS3Key())
			if err != nil {
				l.Errorw("file.get_files_failed.presign_download_error", "err", err)
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			pbVariant, err := fileVariantToProto(fv, &downloadURL, nil, req.IncludeProprietaryInfo)
			if err != nil {
				l.Errorw("file.get_files_failed.map_variant_error", "err", err)
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			pbVariants = append(pbVariants, pbVariant)
		}

		pbFiles = append(pbFiles, fileToProto(f, pbVariants, req.IncludeProprietaryInfo))
	}

	l.Infow("file.get_files.success", "count", len(pbFiles))
	return &pb.GetFilesResponse{Files: pbFiles}, nil
}

func (s *service) DeleteFiles(ctx context.Context, req *pb.DeleteFilesRequest) (*pb.DeleteFilesResponse, error) {
	l := s.log.With("op", "delete_files", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(req.FileIds) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrFileIdsIsRequired.Error())
	}
	if len(req.FileIds) > 1000 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyFiles.Error())
	}

	uow := s.fileDataProvider.newUOW()
	defer uow.Close()

	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("file.delete_files_failed.begin_transaction_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	var files []*file.File
	if len(req.FileIds) == 1 {
		f, err := s.fileDataProvider.getFileByIDLocked(ctx, req.FileIds[0], uow)
		if err != nil {
			l.Errorw("file.delete_files_failed.get_file_locked_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
		if f == nil {
			return nil, grpcstatus.Error(codes.NotFound, ErrFileNotFound.Error())
		}
		files = []*file.File{f}
	} else {
		query := models.NewQueryFilesDal(
			models.FileFilterDal{Ids: req.FileIds},
			1, len(req.FileIds),
		)
		f, err := s.fileDataProvider.getFilesLocked(ctx, query, false, uow)
		if err != nil {
			l.Errorw("file.delete_files_failed.get_files_locked_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
		files = f
	}

	if len(files) == 0 {
		return &pb.DeleteFilesResponse{}, nil
	}

	if err := s.deleteFilesWithVariants(ctx, files, uow, "file.delete_files_failed", l); err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("file.delete_files_failed.commit_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	l.Infow("file.delete_files.success", "count", len(files))
	return &pb.DeleteFilesResponse{}, nil
}

func (s *service) UpdateFiles(ctx context.Context, req *pb.UpdateFilesRequest) (*pb.UpdateFilesResponse, error) {
	l := s.log.With("op", "update_files", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(req.Files) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrFileIdsIsRequired.Error())
	}
	if len(req.Files) > 1000 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyFiles.Error())
	}

	parsed, fileIDs, err := s.validateAndStatUpdateRequest(ctx, req, l)
	if err != nil {
		return nil, err
	}

	uow := s.fileDataProvider.newUOW()
	defer uow.Close()

	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("file.update_files_failed.begin_transaction_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	var lockedFiles []*file.File
	if len(fileIDs) == 1 {
		f, err := s.fileDataProvider.getFileByIDLocked(ctx, fileIDs[0], uow)
		if err != nil {
			l.Errorw("file.update_files_failed.get_file_locked_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
		if f == nil {
			return nil, grpcstatus.Error(codes.NotFound, ErrFileNotFound.Error())
		}
		lockedFiles = []*file.File{f}
	} else {
		query := models.NewQueryFilesDal(models.FileFilterDal{Ids: fileIDs}, 1, len(fileIDs))
		lockedFiles, err = s.fileDataProvider.getFilesLocked(ctx, query, false, uow)
		if err != nil {
			l.Errorw("file.update_files_failed.get_files_locked_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
	}

	variantsQuery := models.NewQueryFileVariantsDal(models.FileVariantFilterDal{FileIds: fileIDs}, nil, nil)
	allVariants, err := s.fileVariantDataProvider.getFileVariantsLocked(ctx, variantsQuery, uow)
	if err != nil {
		l.Errorw("file.update_files_failed.get_variants_locked_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	mutations, fileResults, err := s.buildUpdateMutations(parsed, lockedFiles, allVariants, time.Now(), l)
	if err != nil {
		return nil, err
	}

	if err := s.applyUpdateMutations(ctx, mutations, uow, l); err != nil {
		return nil, err
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("file.update_files_failed.commit_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	pbFiles, err := s.buildUpdateFilesProto(ctx, fileResults, l)
	if err != nil {
		return nil, err
	}

	l.Infow("file.update_files.success", "updated", len(mutations.filesToUpdate), "deleted", len(mutations.filesToDelete))
	return &pb.UpdateFilesResponse{Files: pbFiles}, nil
}

func (s *service) DeleteNotInUseFiles(ctx context.Context, workerID string, batchSize int, uow *uow.UnitOfWork) error {
	l := s.log.With("op", "delete_not_in_use_files", "worker_id", workerID, "batch_size", batchSize)

	needToCommit := false
	if uow == nil {
		uow = s.fileDataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("file.delete_not_in_use_files_failed.begin_transaction_error", "err", err)
			return ErrDeleteNotInUseFiles
		}
	}

	updatedAtBefore := time.Now().Add(-6 * time.Hour)
	createdAtBefore := time.Now().Add(-24 * time.Hour)

	query := models.NewQueryFilesDal(
		models.FileFilterDal{
			Statuses: []int16{
				int16(file.FileStatusWaitingForUpload), int16(file.FileStatusUploaded),
				int16(file.FileStatusProcessing), int16(file.FileStatusReady),
			},
			CreatedTo: &createdAtBefore,
			UpdatedTo: &updatedAtBefore,
		}, 1, batchSize,
	)

	files, err := s.fileDataProvider.getFilesLocked(ctx, query, true, uow)
	if err != nil {
		l.Errorw("file.delete_not_in_use_files_failed.get_files_locked_error", "err", err)
		return ErrDeleteNotInUseFiles
	}
	if len(files) == 0 {
		return nil
	}

	if err := s.deleteFilesWithVariants(ctx, files, uow, "file.delete_not_in_use_files_failed", l); err != nil {
		return ErrDeleteNotInUseFiles
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("file.delete_not_in_use_files_failed.commit_error", "err", err)
			return ErrDeleteNotInUseFiles
		}
	}

	l.Infow("file.delete_not_in_use_files.success", "count", len(files))
	return nil
}

func (s *service) DeleteFilesByIDs(ctx context.Context, workerID string, fileIDs []string, uow *uow.UnitOfWork) error {
	l := s.log.With("op", "delete_files_by_ids", "worker_id", workerID)

	if len(fileIDs) == 0 {
		return nil
	}

	needToCommit := false
	if uow == nil {
		uow = s.fileDataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("file.delete_files_by_ids_failed.begin_transaction_error", "err", err)
			return ErrDeleteFilesByIDs
		}
	}

	query := models.NewQueryFilesDal(models.FileFilterDal{Ids: fileIDs}, 1, len(fileIDs))
	files, err := s.fileDataProvider.getFilesLocked(ctx, query, false, uow)
	if err != nil {
		l.Errorw("file.delete_files_by_ids_failed.get_files_locked_error", "err", err)
		return ErrDeleteFilesByIDs
	}
	if len(files) == 0 {
		return nil
	}

	if err := s.deleteFilesWithVariants(ctx, files, uow, "file.delete_files_by_ids_failed", l); err != nil {
		return ErrDeleteFilesByIDs
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("file.delete_files_by_ids_failed.commit_error", "err", err)
			return ErrDeleteFilesByIDs
		}
	}

	l.Infow("file.delete_files_by_ids.success", "count", len(files))
	return nil
}

func (s *service) deleteFilesWithVariants(
	ctx context.Context,
	files []*file.File,
	uow *uow.UnitOfWork,
	logPrefix string,
	l *zap.SugaredLogger,
) error {
	fileIDs := make([]string, len(files))
	for i, f := range files {
		fileIDs[i] = f.GetID()
	}

	variantsQuery := models.NewQueryFileVariantsDal(
		models.FileVariantFilterDal{FileIds: fileIDs},
		nil, nil,
	)
	variants, err := s.fileVariantDataProvider.getFileVariantsLocked(ctx, variantsQuery, uow)
	if err != nil {
		l.Errorw(logPrefix+".get_variants_locked_error", "err", err)
		return err
	}

	s3KeysByBucket := make(map[filevariant.BucketType][]string)
	for _, fv := range variants {
		s3KeysByBucket[fv.GetBucket()] = append(s3KeysByBucket[fv.GetBucket()], fv.GetS3Key())
	}

	if err := s.fileVariantDataProvider.deleteFileVariants(ctx, variants, uow); err != nil {
		l.Errorw(logPrefix+".delete_variants_error", "err", err)
		return err
	}

	if err := s.fileDataProvider.deleteFiles(ctx, files, uow); err != nil {
		l.Errorw(logPrefix+".delete_files_error", "err", err)
		return err
	}

	tasks := make([]filedeletionfroms3tasksservice.FileDeletionFromS3Task, 0, len(s3KeysByBucket))
	for bucket, keys := range s3KeysByBucket {
		tasks = append(tasks, filedeletionfroms3tasksservice.FileDeletionFromS3Task{
			Bucket: int16(bucket),
			S3Keys: keys,
		})
	}
	if err := s.fileDeletionFromS3TasksService.CreateFileDeletionFromS3Tasks(ctx, "", tasks, uow); err != nil {
		l.Errorw(logPrefix+".create_deletion_tasks_error", "err", err)
		return err
	}

	return nil
}

func (s *service) validateAndStatUpdateRequest(
	ctx context.Context,
	req *pb.UpdateFilesRequest,
	l *zap.SugaredLogger,
) ([]fileParsed, []string, error) {
	parsed, fileIDs, verr := s.parseUpdateRequest(req)
	if len(verr) > 0 {
		return nil, nil, verr.ToStatus()
	}
	statVerr, err := s.statAddVariants(ctx, parsed, l)
	if err != nil {
		return nil, nil, err
	}
	if len(statVerr) > 0 {
		return nil, nil, statVerr.ToStatus()
	}
	return parsed, fileIDs, nil
}

func (s *service) parseUpdateRequest(req *pb.UpdateFilesRequest) ([]fileParsed, []string, validationerror.ValidationError) {
	parsed := make([]fileParsed, len(req.Files))
	fileIDs := make([]string, 0, len(req.Files))
	verr := make(validationerror.ValidationError)

	for i, fr := range req.Files {
		if fr.Id == "" {
			verr[fmt.Sprintf("files[%d].id", i)] = ErrFileIdsIsRequired.Error()
			continue
		}
		fp := fileParsed{id: fr.Id}
		fileIDs = append(fileIDs, fr.Id)

		needPurpose := len(fr.AddVariants) > 0
		if !needPurpose {
			for _, uv := range fr.UpdateVariants {
				if uv.VariantType != nil || uv.Meta != nil {
					needPurpose = true
					break
				}
			}
		}

		if fr.FilePurpose != nil {
			purpose, err := file.FilePurposeFromString(*fr.FilePurpose)
			if err != nil {
				verr[fmt.Sprintf("files[%d].file_purpose", i)] = err.Error()
			} else {
				fp.purpose = &purpose
			}
		} else if needPurpose {
			verr[fmt.Sprintf("files[%d].file_purpose", i)] = ErrFilePurposeIsRequired.Error()
		}

		if fr.Status != nil {
			status, err := file.FileStatusFromString(*fr.Status)
			if err != nil {
				verr[fmt.Sprintf("files[%d].status", i)] = err.Error()
			} else {
				fp.newStatus = &status
			}
		}

		if len(fr.DeleteVariants) > 0 {
			fp.deleteVariantIDs = make(map[string]struct{}, len(fr.DeleteVariants))
			for _, id := range fr.DeleteVariants {
				fp.deleteVariantIDs[id] = struct{}{}
			}
		}

		fp.updateVariants = make([]variantUpdateParsed, 0, len(fr.UpdateVariants))
		for j, uv := range fr.UpdateVariants {
			if uv.Id == "" {
				verr[fmt.Sprintf("files[%d].update_variants[%d].id", i, j)] = ErrVariantIdsIsRequired.Error()
				continue
			}
			uvp := variantUpdateParsed{id: uv.Id}
			if uv.VariantType != nil && fp.purpose != nil {
				vt, err := filevariant.VariantTypeFromString(*uv.VariantType)
				if err != nil {
					verr[fmt.Sprintf("files[%d].update_variants[%d].variant_type", i, j)] = err.Error()
				} else if err := filevariant.ValidateVariantType(vt, *fp.purpose); err != nil {
					verr[fmt.Sprintf("files[%d].update_variants[%d].variant_type", i, j)] = err.Error()
				} else {
					uvp.variantType = &vt
				}
			}
			if uv.Meta != nil && fp.purpose != nil {
				var meta filevariant.Meta
				if err := json.Unmarshal([]byte(*uv.Meta), &meta); err != nil {
					verr[fmt.Sprintf("files[%d].update_variants[%d].meta", i, j)] = filevariant.ErrUmarshalMeta.Error()
				} else if err := filevariant.ValidateMeta(meta, *fp.purpose); err != nil {
					verr[fmt.Sprintf("files[%d].update_variants[%d].meta", i, j)] = err.Error()
				} else {
					uvp.meta = &meta
				}
			}
			fp.updateVariants = append(fp.updateVariants, uvp)
		}

		fp.addVariants = make([]variantAddParsed, 0, len(fr.AddVariants))
		for j, av := range fr.AddVariants {
			avp := variantAddParsed{s3Key: av.S3Key}
			if av.S3Key == "" {
				verr[fmt.Sprintf("files[%d].add_variants[%d].s3_key", i, j)] = ErrS3KeyIsRequired.Error()
			}
			if fp.purpose != nil {
				bucket, _ := s.purposeStorage(*fp.purpose)
				avp.bucket = bucket
				vt, err := filevariant.VariantTypeFromString(av.VariantType)
				if err != nil {
					verr[fmt.Sprintf("files[%d].add_variants[%d].variant_type", i, j)] = err.Error()
				} else if err := filevariant.ValidateVariantType(vt, *fp.purpose); err != nil {
					verr[fmt.Sprintf("files[%d].add_variants[%d].variant_type", i, j)] = err.Error()
				} else {
					avp.variantType = vt
				}
				if av.Meta != "" {
					var meta filevariant.Meta
					if err := json.Unmarshal([]byte(av.Meta), &meta); err != nil {
						verr[fmt.Sprintf("files[%d].add_variants[%d].meta", i, j)] = filevariant.ErrUmarshalMeta.Error()
					} else if err := filevariant.ValidateMeta(meta, *fp.purpose); err != nil {
						verr[fmt.Sprintf("files[%d].add_variants[%d].meta", i, j)] = err.Error()
					} else {
						avp.meta = meta
					}
				}
			}
			fp.addVariants = append(fp.addVariants, avp)
		}

		parsed[i] = fp
	}

	return parsed, fileIDs, verr
}

func (s *service) statAddVariants(
	ctx context.Context,
	parsed []fileParsed,
	l *zap.SugaredLogger,
) (validationerror.ValidationError, error) {
	verr := make(validationerror.ValidationError)
	for i := range parsed {
		for j := range parsed[i].addVariants {
			avp := &parsed[i].addVariants[j]
			if avp.s3Key == "" {
				continue
			}
			info, err := s.bucketStat(ctx, avp.bucket, avp.s3Key)
			if err != nil {
				l.Errorw("file.update_files_failed.s3_stat_error", "err", err, "key", avp.s3Key)
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			if info == nil {
				verr[fmt.Sprintf("files[%d].add_variants[%d].s3_key", i, j)] = ErrS3ObjectNotFound.Error()
				continue
			}
			avp.mimeType = info.ContentType
			avp.size = info.Size
		}
	}
	return verr, nil
}

func (s *service) buildUpdateMutations(
	parsed []fileParsed,
	lockedFiles []*file.File,
	allVariants []*filevariant.FileVariant,
	now time.Time,
	l *zap.SugaredLogger,
) (*updateMutations, []fileResult, error) {
	fileMap := make(map[string]*file.File, len(lockedFiles))
	for _, f := range lockedFiles {
		fileMap[f.GetID()] = f
	}

	variantsByFileID := make(map[string][]*filevariant.FileVariant)
	variantByID := make(map[string]*filevariant.FileVariant)
	for _, fv := range allVariants {
		variantsByFileID[fv.GetFileID()] = append(variantsByFileID[fv.GetFileID()], fv)
		variantByID[fv.GetID()] = fv
	}

	mut := &updateMutations{
		s3KeysByBucket: make(map[filevariant.BucketType][]string),
	}
	fileResults := make([]fileResult, 0, len(parsed))

	for _, fp := range parsed {
		f, exists := fileMap[fp.id]
		if !exists {
			continue
		}

		if fp.newStatus != nil {
			if err := f.SetStatus(*fp.newStatus); err != nil {
				return nil, nil, grpcstatus.Errorf(codes.InvalidArgument, "file %s: %s", fp.id, err.Error())
			}
		}

		fileVariants := variantsByFileID[fp.id]

		for dvID := range fp.deleteVariantIDs {
			fv, ok := variantByID[dvID]
			if !ok || fv.GetFileID() != fp.id {
				return nil, nil, grpcstatus.Errorf(codes.InvalidArgument, "file %s: variant %s not found", fp.id, dvID)
			}
			mut.variantsToDelete = append(mut.variantsToDelete, fv)
			mut.s3KeysByBucket[fv.GetBucket()] = append(mut.s3KeysByBucket[fv.GetBucket()], fv.GetS3Key())
		}

		for _, uvp := range fp.updateVariants {
			fv, ok := variantByID[uvp.id]
			if !ok || fv.GetFileID() != fp.id {
				return nil, nil, grpcstatus.Errorf(codes.InvalidArgument, "file %s: variant %s not found", fp.id, uvp.id)
			}
			if _, deleted := fp.deleteVariantIDs[uvp.id]; deleted {
				return nil, nil, grpcstatus.Errorf(codes.InvalidArgument, "file %s: variant %s is in both update_variants and delete_variants", fp.id, uvp.id)
			}
			if uvp.variantType != nil {
				if err := fv.SetVariant(*uvp.variantType, *fp.purpose); err != nil {
					return nil, nil, grpcstatus.Errorf(codes.InvalidArgument, "file %s: variant %s: %s", fp.id, uvp.id, err.Error())
				}
			}
			if uvp.meta != nil {
				if err := fv.SetMeta(*uvp.meta, *fp.purpose); err != nil {
					return nil, nil, grpcstatus.Errorf(codes.InvalidArgument, "file %s: variant %s: %s", fp.id, uvp.id, err.Error())
				}
			}
			fv.SetUpdatedAt(&now)
			mut.variantsToUpdate = append(mut.variantsToUpdate, fv)
		}

		newVariants := make([]*filevariant.FileVariant, 0, len(fp.addVariants))
		for _, avp := range fp.addVariants {
			metaJSON, _ := json.Marshal(avp.meta)
			fv, err := filevariant.FromStorage("", fp.id, int16(avp.variantType), int16(avp.bucket), avp.s3Key, avp.mimeType, avp.size, metaJSON, now, now)
			if err != nil {
				l.Errorw("file.update_files_failed.build_variant_error", "err", err)
				return nil, nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			newVariants = append(newVariants, fv)
			mut.variantsToCreate = append(mut.variantsToCreate, fv)
		}

		finalVariants := make([]*filevariant.FileVariant, 0, len(fileVariants)+len(newVariants))
		for _, fv := range fileVariants {
			if _, deleted := fp.deleteVariantIDs[fv.GetID()]; !deleted {
				finalVariants = append(finalVariants, fv)
			}
		}
		finalVariants = append(finalVariants, newVariants...)

		if len(finalVariants) == 0 {
			mut.filesToDelete = append(mut.filesToDelete, f)
		} else {
			f.SetUpdatedAt(&now)
			mut.filesToUpdate = append(mut.filesToUpdate, f)
			fileResults = append(fileResults, fileResult{f: f, finalVariants: finalVariants})
		}
	}

	return mut, fileResults, nil
}

func (s *service) applyUpdateMutations(
	ctx context.Context,
	mut *updateMutations,
	uow *uow.UnitOfWork,
	l *zap.SugaredLogger,
) error {
	if err := s.fileVariantDataProvider.deleteFileVariants(ctx, mut.variantsToDelete, uow); err != nil {
		l.Errorw("file.update_files_failed.delete_variants_error", "err", err)
		return grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if err := s.fileVariantDataProvider.updateFileVariants(ctx, mut.variantsToUpdate, uow); err != nil {
		l.Errorw("file.update_files_failed.update_variants_error", "err", err)
		return grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if err := s.fileVariantDataProvider.createFileVariants(ctx, mut.variantsToCreate, uow); err != nil {
		l.Errorw("file.update_files_failed.create_variants_error", "err", err)
		return grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if err := s.fileDataProvider.deleteFiles(ctx, mut.filesToDelete, uow); err != nil {
		l.Errorw("file.update_files_failed.delete_files_error", "err", err)
		return grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if err := s.fileDataProvider.updateFiles(ctx, mut.filesToUpdate, uow); err != nil {
		l.Errorw("file.update_files_failed.update_files_error", "err", err)
		return grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	s3Tasks := make([]filedeletionfroms3tasksservice.FileDeletionFromS3Task, 0, len(mut.s3KeysByBucket))
	for bucket, keys := range mut.s3KeysByBucket {
		s3Tasks = append(s3Tasks, filedeletionfroms3tasksservice.FileDeletionFromS3Task{
			Bucket: int16(bucket),
			S3Keys: keys,
		})
	}
	if err := s.fileDeletionFromS3TasksService.CreateFileDeletionFromS3Tasks(ctx, "", s3Tasks, uow); err != nil {
		l.Errorw("file.update_files_failed.create_deletion_tasks_error", "err", err)
		return grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	return nil
}

func (s *service) buildUpdateFilesProto(
	ctx context.Context,
	results []fileResult,
	l *zap.SugaredLogger,
) ([]*pb.File, error) {
	pbFiles := make([]*pb.File, 0, len(results))
	for _, res := range results {
		pbVariants := make([]*pb.FileVariant, 0, len(res.finalVariants))
		for _, fv := range res.finalVariants {
			downloadURL, err := s.bucketPresignDownload(ctx, fv.GetBucket(), fv.GetS3Key())
			if err != nil {
				l.Errorw("file.update_files_failed.presign_download_error", "err", err)
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			pbVariant, err := fileVariantToProto(fv, &downloadURL, nil, true)
			if err != nil {
				l.Errorw("file.update_files_failed.map_variant_error", "err", err)
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			pbVariants = append(pbVariants, pbVariant)
		}
		pbFiles = append(pbFiles, fileToProto(res.f, pbVariants, true))
	}
	return pbFiles, nil
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

func (s *service) bucketStat(ctx context.Context, bucket filevariant.BucketType, key string) (*s3transport.ObjectInfo, error) {
	if bucket == filevariant.UserPhotosBucket {
		return s.fileStorage.userPhotoStat(ctx, key)
	}
	return s.fileStorage.chatFileStat(ctx, key)
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
