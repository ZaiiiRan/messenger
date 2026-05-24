package file

type FilePurpose int16

const (
	PurposeUserPhoto  FilePurpose = 0
	PurposeChatPhoto  FilePurpose = 1
	PurposeChatVideo  FilePurpose = 2
	PurposeChatAudio  FilePurpose = 3
	PurposeChatFile   FilePurpose = 4
	PurposeChatAvatar FilePurpose = 5
)

func (p FilePurpose) String() string {
	switch p {
	case PurposeUserPhoto:
		return "user_photo"
	case PurposeChatPhoto:
		return "chat_photo"
	case PurposeChatVideo:
		return "chat_video"
	case PurposeChatAudio:
		return "chat_audio"
	case PurposeChatFile:
		return "chat_file"
	case PurposeChatAvatar:
		return "chat_avatar"
	default:
		return "unknown"
	}
}

func FilePurposeFromString(s string) (FilePurpose, error) {
	switch s {
	case "user_photo":
		return PurposeUserPhoto, nil
	case "chat_photo":
		return PurposeChatPhoto, nil
	case "chat_video":
		return PurposeChatVideo, nil
	case "chat_audio":
		return PurposeChatAudio, nil
	case "chat_file":
		return PurposeChatFile, nil
	case "chat_avatar":
		return PurposeChatAvatar, nil
	default:
		return 0, ErrUnknownFilePurpose
	}
}
