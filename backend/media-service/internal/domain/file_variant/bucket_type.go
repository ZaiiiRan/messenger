package filevariant

type BucketType int16

const (
	UserPhotosBucket BucketType = 0
	ChatFilesBucket  BucketType = 1
)

func (b BucketType) String() string {
	switch b {
	case UserPhotosBucket:
		return "user_photos"
	case ChatFilesBucket:
		return "chat_files"
	}
	return "unknown"
}

func BucketTypeFromString(s string) (BucketType, bool) {
	switch s {
	case "user_photos":
		return UserPhotosBucket, true
	case "chat_files":
		return ChatFilesBucket, true
	}
	return 0, false
}
