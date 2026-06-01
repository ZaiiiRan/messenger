package processing

import (
	file "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file"
	filevariant "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file_variant"
)

type ProcessingType int16

const (
	ProcessingUserPhoto  ProcessingType = 0
	ProcessingChatAvatar ProcessingType = 1
	ProcessingChatPhoto  ProcessingType = 2
	ProcessingChatVideo  ProcessingType = 3
	ProcessingChatAudio  ProcessingType = 4
)

func (p ProcessingType) String() string {
	switch p {
	case ProcessingUserPhoto:
		return "user_photo"
	case ProcessingChatAvatar:
		return "chat_avatar"
	case ProcessingChatPhoto:
		return "chat_photo"
	case ProcessingChatVideo:
		return "chat_video"
	case ProcessingChatAudio:
		return "chat_audio"
	}
	return "unknown"
}

func ProcessingTypeFromString(s string) (ProcessingType, error) {
	switch s {
	case "user_photo":
		return ProcessingUserPhoto, nil
	case "chat_avatar":
		return ProcessingChatAvatar, nil
	case "chat_photo":
		return ProcessingChatPhoto, nil
	case "chat_video":
		return ProcessingChatVideo, nil
	case "chat_audio":
		return ProcessingChatAudio, nil
	default:
		return 0, ErrUnknownProcessingType
	}
}

func (p ProcessingType) Purpose() file.FilePurpose {
	switch p {
	case ProcessingUserPhoto:
		return file.PurposeUserPhoto
	case ProcessingChatAvatar:
		return file.PurposeChatAvatar
	case ProcessingChatPhoto:
		return file.PurposeChatPhoto
	case ProcessingChatVideo:
		return file.PurposeChatVideo
	case ProcessingChatAudio:
		return file.PurposeChatAudio
	}
	return file.PurposeChatFile
}

func (p ProcessingType) Category() MediaCategory {
	switch p {
	case ProcessingUserPhoto, ProcessingChatAvatar, ProcessingChatPhoto:
		return CategoryImage
	case ProcessingChatVideo:
		return CategoryVideo
	case ProcessingChatAudio:
		return CategoryAudio
	}
	return CategoryImage
}

func (p ProcessingType) Bucket() filevariant.BucketType {
	if p == ProcessingUserPhoto {
		return filevariant.UserPhotosBucket
	}
	return filevariant.ChatFilesBucket
}

func (p ProcessingType) RequiresCrop() bool {
	return p == ProcessingUserPhoto || p == ProcessingChatAvatar
}
