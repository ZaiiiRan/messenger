package file

type FileStatus int16

const (
	FileStatusWaitingForUpload FileStatus = 0
	FileStatusUploaded         FileStatus = 1
	FileStatusProcessing       FileStatus = 2
	FileStatusReady            FileStatus = 3
	FileStatusInUse            FileStatus = 4
)

func (s FileStatus) String() string {
	switch s {
	case FileStatusWaitingForUpload:
		return "waiting_for_upload"
	case FileStatusUploaded:
		return "uploaded"
	case FileStatusProcessing:
		return "processing"
	case FileStatusReady:
		return "ready"
	case FileStatusInUse:
		return "in_use"
	}
	return "unknown"
}

func FileStatusFromString(s string) (FileStatus, error) {
	switch s {
	case "waiting_for_upload":
		return FileStatusWaitingForUpload, nil
	case "uploaded":
		return FileStatusUploaded, nil
	case "processing":
		return FileStatusProcessing, nil
	case "ready":
		return FileStatusReady, nil
	case "in_use":
		return FileStatusInUse, nil
	default:
		return 0, ErrUnknownFileStatus
	}
}
