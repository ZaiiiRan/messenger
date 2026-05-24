package filevariant

type VariantType int16

const (
	VariantOriginal     VariantType = 0
	VariantAvatar       VariantType = 1
	VariantAvatarSmall  VariantType = 2
	VariantAvatarMiddle VariantType = 3
	VariantCompressed   VariantType = 4
	Variant360p         VariantType = 5
	Variant720p         VariantType = 6
	Variant1080p        VariantType = 7
	VariantThumb        VariantType = 8
)

func (v VariantType) String() string {
	switch v {
	case VariantOriginal:
		return "original"
	case VariantAvatar:
		return "avatar"
	case VariantAvatarSmall:
		return "avatar_small"
	case VariantAvatarMiddle:
		return "avatar_middle"
	case VariantCompressed:
		return "compressed"
	case Variant360p:
		return "360p"
	case Variant720p:
		return "720p"
	case Variant1080p:
		return "1080p"
	case VariantThumb:
		return "thumb"
	default:
		return "unknown"
	}
}

func VariantTypeFromString(s string) (VariantType, error) {
	switch s {
	case "original":
		return VariantOriginal, nil
	case "avatar":
		return VariantAvatar, nil
	case "avatar_small":
		return VariantAvatarSmall, nil
	case "avatar_middle":
		return VariantAvatarMiddle, nil
	case "compressed":
		return VariantCompressed, nil
	case "360p":
		return Variant360p, nil
	case "720p":
		return Variant720p, nil
	case "1080p":
		return Variant1080p, nil
	case "thumb":
		return VariantThumb, nil
	default:
		return 0, ErrUnknownVariantType
	}
}
