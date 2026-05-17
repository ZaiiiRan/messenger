package models

import (
	"encoding/json"

	privacysettings "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/privacy_settings"
)

type V1PrivacySettingsDal struct {
	Id       int64           `db:"id" json:"id"`
	UserId   string          `db:"user_id" json:"user_id"`
	Settings json.RawMessage `db:"settings" json:"settings"`
}

type V1PrivacySettingsAggregateDal struct {
	Avatar           V1PrivacySettingDal `json:"avatar"`
	Photos           V1PrivacySettingDal `json:"photos"`
	PhoneNumber      V1PrivacySettingDal `json:"phone_number"`
	Email            V1PrivacySettingDal `json:"email"`
	Birthdate        V1PrivacySettingDal `json:"birthdate"`
	OnlineStatus     V1PrivacySettingDal `json:"online_status"`
	FirstDialogsInit V1PrivacySettingDal `json:"first_dialogs_init"`
	GroupChatInvites V1PrivacySettingDal `json:"group_chat_invites"`
}

type V1PrivacySettingDal struct {
	Value      int16    `json:"value"`
	Favourites []string `json:"favourites"`
	Exceptions []string `json:"exceptions"`
}

func V1PrivacySettingsDalFromDomain(userId string, s *privacysettings.PrivacySettings) (V1PrivacySettingsDal, error) {
	if s == nil {
		return V1PrivacySettingsDal{}, nil
	}

	privacySettingsAggregateDal := V1PrivacySettingsAggregateDal{
		Avatar: V1PrivacySettingDal{
			Value:      int16(s.GetAvatar().GetValue()),
			Favourites: s.GetAvatar().GetFavourites(),
			Exceptions: s.GetAvatar().GetExceptions(),
		},
		Photos: V1PrivacySettingDal{
			Value:      int16(s.GetPhotos().GetValue()),
			Favourites: s.GetPhotos().GetFavourites(),
			Exceptions: s.GetPhotos().GetExceptions(),
		},
		PhoneNumber: V1PrivacySettingDal{
			Value:      int16(s.GetPhoneNumber().GetValue()),
			Favourites: s.GetPhoneNumber().GetFavourites(),
			Exceptions: s.GetPhoneNumber().GetExceptions(),
		},
		Email: V1PrivacySettingDal{
			Value:      int16(s.GetEmail().GetValue()),
			Favourites: s.GetEmail().GetFavourites(),
			Exceptions: s.GetEmail().GetExceptions(),
		},
		Birthdate: V1PrivacySettingDal{
			Value:      int16(s.GetBirthdate().GetValue()),
			Favourites: s.GetBirthdate().GetFavourites(),
			Exceptions: s.GetBirthdate().GetExceptions(),
		},
		OnlineStatus: V1PrivacySettingDal{
			Value:      int16(s.GetOnlineStatus().GetValue()),
			Favourites: s.GetOnlineStatus().GetFavourites(),
			Exceptions: s.GetOnlineStatus().GetExceptions(),
		},
		FirstDialogsInit: V1PrivacySettingDal{
			Value:      int16(s.GetFirstDialogsInit().GetValue()),
			Favourites: s.GetFirstDialogsInit().GetFavourites(),
			Exceptions: s.GetFirstDialogsInit().GetExceptions(),
		},
		GroupChatInvites: V1PrivacySettingDal{
			Value:      int16(s.GetGroupChatInvites().GetValue()),
			Favourites: s.GetGroupChatInvites().GetFavourites(),
			Exceptions: s.GetGroupChatInvites().GetExceptions(),
		},
	}

	settingsJson, err := json.Marshal(privacySettingsAggregateDal)
	if err != nil {
		return V1PrivacySettingsDal{}, err
	}

	return V1PrivacySettingsDal{
		UserId:   userId,
		Settings: settingsJson,
	}, nil
}

func (s V1PrivacySettingsDal) IsNull() bool { return false }
func (s V1PrivacySettingsDal) Index(i int) any {
	switch i {
	case 0:
		return s.Id
	case 1:
		return s.UserId
	case 2:
		return s.Settings
	default:
		return nil
	}
}

func (s V1PrivacySettingsDal) ToDomain() (*privacysettings.PrivacySettings, error) {
	var settings V1PrivacySettingsAggregateDal
	err := json.Unmarshal(s.Settings, &settings)
	if err != nil {
		return nil, err
	}

	avatar := privacysettings.PrivacySettingFromStorage(
		settings.Avatar.Value,
		settings.Avatar.Favourites,
		settings.Avatar.Exceptions,
	)
	photos := privacysettings.PrivacySettingFromStorage(
		settings.Photos.Value,
		settings.Photos.Favourites,
		settings.Photos.Exceptions,
	)
	phoneNumber := privacysettings.PrivacySettingFromStorage(
		settings.PhoneNumber.Value,
		settings.PhoneNumber.Favourites,
		settings.PhoneNumber.Exceptions,
	)
	email := privacysettings.PrivacySettingFromStorage(
		settings.Email.Value,
		settings.Email.Favourites,
		settings.Email.Exceptions,
	)
	birthdate := privacysettings.PrivacySettingFromStorage(
		settings.Birthdate.Value,
		settings.Birthdate.Favourites,
		settings.Birthdate.Exceptions,
	)
	onlineStatus := privacysettings.PrivacySettingFromStorage(
		settings.OnlineStatus.Value,
		settings.OnlineStatus.Favourites,
		settings.OnlineStatus.Exceptions,
	)
	firstDialogsInit := privacysettings.PrivacySettingFromStorage(
		settings.FirstDialogsInit.Value,
		settings.FirstDialogsInit.Favourites,
		settings.FirstDialogsInit.Exceptions,
	)
	groupChatInvites := privacysettings.PrivacySettingFromStorage(
		settings.GroupChatInvites.Value,
		settings.GroupChatInvites.Favourites,
		settings.GroupChatInvites.Exceptions,
	)

	return privacysettings.FromStorage(
		avatar,
		photos,
		phoneNumber,
		email,
		birthdate,
		onlineStatus,
		firstDialogsInit,
		groupChatInvites,
	), nil
}
