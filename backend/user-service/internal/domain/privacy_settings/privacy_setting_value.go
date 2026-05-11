package privacysettings

type PrivacyValue int16

var (
	All     PrivacyValue = 0
	Friends PrivacyValue = 1
	None    PrivacyValue = 2
)

func (p PrivacyValue) String() string {
	switch p {
	case All:
		return "all"
	case Friends:
		return "friends"
	case None:
		return "none"
	default:
		return ""
	}
}

func ToPrivacyValue(value string) (PrivacyValue, error) {
	switch value {
	case "all":
		return All, nil
	case "friends":
		return Friends, nil
	case "none":
		return None, nil
	default:
		return All, ErrUnknownPrivacyValue
	}
}
