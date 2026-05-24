package utils

import (
	"strings"

	pb "github.com/ZaiiiRan/messenger/backend/file-service/gen/go/file/v1"
)

func SanitizeCreateFilesByUserRequest(req *pb.CreateFilesByUserRequest) {
	if req == nil {
		return
	}
	if len(req.Files) == 0 {
		return
	}
	for _, file := range req.Files {
		sanitizeCreateFileRequest(file)
	}
}

func sanitizeCreateFileRequest(req *pb.CreateFileRequest) {
	if req == nil {
		return
	}
	req.MimeType = strings.ToLower(strings.TrimSpace(req.MimeType))
	if req.Size < 0 {
		req.Size = 0
	}
	req.FilePurpose = strings.ToLower(strings.TrimSpace(req.FilePurpose))
}

func SanitizeConfirmUploadsByUserRequest(req *pb.ConfirmUploadsByUserRequest) {
	if req == nil {
		return
	}
	req.FileIds = sanitizeUniqueArray(sanitizeStringArray(req.FileIds))
}

func SanitizeGetFilesByUserRequest(req *pb.GetFilesByUserRequest) {
	if req == nil {
		return
	}
	req.FileIds = sanitizeUniqueArray(sanitizeStringArray(req.FileIds))
}

func SanitizeGetFilesRequest(req *pb.GetFilesRequest) {
	if req == nil {
		return
	}
	req.FileIds = sanitizeUniqueArray(sanitizeStringArray(req.FileIds))
}

func SanitizeDeleteFilesRequest(req *pb.DeleteFilesRequest) {
	if req == nil {
		return
	}
	req.FileIds = sanitizeUniqueArray(sanitizeStringArray(req.FileIds))
}

func SanitizeUpdateFilesRequest(req *pb.UpdateFilesRequest) {
	if req == nil {
		return
	}
	for _, file := range req.Files {
		sanitizeUpdateFileRequest(file)
	}
}

func sanitizeUpdateFileRequest(req *pb.UpdateFileRequest) {
	if req == nil {
		return
	}
	req.Id = strings.TrimSpace(req.Id)
	req.Status = sanitizeLowerStringPtr(req.Status)
	req.FilePurpose = sanitizeLowerStringPtr(req.FilePurpose)

	updateVariantsById := make(map[string]*pb.UpdateVariantRequest)
	for _, variant := range req.UpdateVariants {
		sanitizeUpdateVariantRequest(variant)
		if variant.Id != "" {
			updateVariantsById[variant.Id] = variant
		}
	}

	req.UpdateVariants = make([]*pb.UpdateVariantRequest, 0, len(updateVariantsById))
	for _, variant := range updateVariantsById {
		req.UpdateVariants = append(req.UpdateVariants, variant)
	}

	addVariantsByS3Key := make(map[string]*pb.AddVariantRequest)
	for _, variant := range req.AddVariants {
		sanitizeAddVariantRequest(variant)
		if variant.S3Key != "" {
			addVariantsByS3Key[variant.S3Key] = variant
		}
	}

	req.AddVariants = make([]*pb.AddVariantRequest, 0, len(addVariantsByS3Key))
	for _, variant := range addVariantsByS3Key {
		req.AddVariants = append(req.AddVariants, variant)
	}

	req.DeleteVariants = sanitizeUniqueArray(sanitizeStringArray(req.DeleteVariants))
}

func sanitizeUpdateVariantRequest(req *pb.UpdateVariantRequest) {
	if req == nil {
		return
	}
	req.Id = strings.TrimSpace(req.Id)
	req.VariantType = sanitizeLowerStringPtr(req.VariantType)
	req.Meta = sanitizeStringPtr(req.Meta)
}

func sanitizeAddVariantRequest(req *pb.AddVariantRequest) {
	if req == nil {
		return
	}
	req.VariantType = strings.ToLower(strings.TrimSpace(req.VariantType))
	req.S3Key = strings.TrimSpace(req.S3Key)
	req.Meta = strings.TrimSpace(req.Meta)
}

func sanitizeUniqueArray[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	out := s[:0]
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func sanitizeStringArray(arr []string) []string {
	var sanitized []string
	for _, str := range arr {
		str = strings.TrimSpace(str)
		if str != "" {
			sanitized = append(sanitized, str)
		}
	}
	return sanitized
}

func sanitizeLowerStringArray(arr []string) []string {
	var sanitized []string
	for _, str := range arr {
		str = strings.ToLower(strings.TrimSpace(str))
		if str != "" {
			sanitized = append(sanitized, str)
		}
	}
	return sanitized
}

func sanitizeStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func sanitizeLowerStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.ToLower(strings.TrimSpace(*s))
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
