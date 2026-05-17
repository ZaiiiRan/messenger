package privacysettings

type PrivacySetting struct {
	value      int16
	favourites []string
	exceptions []string
}

func PrivacySettingFromStorage(
	value int16,
	favourites []string,
	exceptions []string,
) PrivacySetting {
	return PrivacySetting{
		value:      value,
		favourites: favourites,
		exceptions: exceptions,
	}
}

func (p PrivacySetting) GetValue() PrivacyValue  { return PrivacyValue(p.value) }
func (p PrivacySetting) GetFavourites() []string { return p.favourites }
func (p PrivacySetting) GetExceptions() []string { return p.exceptions }

func (p PrivacySetting) WithValue(v PrivacyValue) PrivacySetting {
	p.value = int16(v)
	return p
}

func (p PrivacySetting) WithFavourites(ids []string) PrivacySetting {
	p.favourites = ids
	return p
}

func (p PrivacySetting) WithExceptions(ids []string) PrivacySetting {
	p.exceptions = ids
	return p
}

type PrivacySettings struct {
	avatar           PrivacySetting
	photos           PrivacySetting
	phoneNumber      PrivacySetting
	email            PrivacySetting
	birthdate        PrivacySetting
	onlineStatus     PrivacySetting
	firstDialogsInit PrivacySetting
	groupChatInvites PrivacySetting
}

func New() *PrivacySettings {
	avatar := PrivacySetting{
		value:      int16(All),
		favourites: make([]string, 0),
		exceptions: make([]string, 0),
	}
	photos := PrivacySetting{
		value:      int16(All),
		favourites: make([]string, 0),
		exceptions: make([]string, 0),
	}
	phoneNumber := PrivacySetting{
		value:      int16(Friends),
		favourites: make([]string, 0),
		exceptions: make([]string, 0),
	}
	email := PrivacySetting{
		value:      int16(Friends),
		favourites: make([]string, 0),
		exceptions: make([]string, 0),
	}
	birthdate := PrivacySetting{
		value:      int16(Friends),
		favourites: make([]string, 0),
		exceptions: make([]string, 0),
	}
	onlineStatus := PrivacySetting{
		value:      int16(Friends),
		favourites: make([]string, 0),
		exceptions: make([]string, 0),
	}
	firstDialogsInit := PrivacySetting{
		value:      int16(All),
		favourites: make([]string, 0),
		exceptions: make([]string, 0),
	}
	groupChatInvites := PrivacySetting{
		value:      int16(Friends),
		favourites: make([]string, 0),
		exceptions: make([]string, 0),
	}
	return &PrivacySettings{
		avatar:           avatar,
		photos:           photos,
		phoneNumber:      phoneNumber,
		email:            email,
		birthdate:        birthdate,
		onlineStatus:     onlineStatus,
		firstDialogsInit: firstDialogsInit,
		groupChatInvites: groupChatInvites,
	}
}

func FromStorage(
	avatar PrivacySetting,
	photos PrivacySetting,
	phoneNumber PrivacySetting,
	email PrivacySetting,
	birthdate PrivacySetting,
	onlineStatus PrivacySetting,
	firstDialogsInit PrivacySetting,
	groupChatInvites PrivacySetting,
) *PrivacySettings {
	return &PrivacySettings{
		avatar:           avatar,
		photos:           photos,
		phoneNumber:      phoneNumber,
		email:            email,
		birthdate:        birthdate,
		onlineStatus:     onlineStatus,
		firstDialogsInit: firstDialogsInit,
		groupChatInvites: groupChatInvites,
	}
}

func (p *PrivacySettings) GetAvatar() PrivacySetting           { return p.avatar }
func (p *PrivacySettings) GetPhotos() PrivacySetting           { return p.photos }
func (p *PrivacySettings) GetPhoneNumber() PrivacySetting      { return p.phoneNumber }
func (p *PrivacySettings) GetEmail() PrivacySetting            { return p.email }
func (p *PrivacySettings) GetBirthdate() PrivacySetting        { return p.birthdate }
func (p *PrivacySettings) GetOnlineStatus() PrivacySetting     { return p.onlineStatus }
func (p *PrivacySettings) GetFirstDialogsInit() PrivacySetting { return p.firstDialogsInit }
func (p *PrivacySettings) GetGroupChatInvites() PrivacySetting { return p.groupChatInvites }

func (p *PrivacySettings) SetAvatar(s PrivacySetting)           { p.avatar = s }
func (p *PrivacySettings) SetPhotos(s PrivacySetting)           { p.photos = s }
func (p *PrivacySettings) SetPhoneNumber(s PrivacySetting)      { p.phoneNumber = s }
func (p *PrivacySettings) SetEmail(s PrivacySetting)            { p.email = s }
func (p *PrivacySettings) SetBirthdate(s PrivacySetting)        { p.birthdate = s }
func (p *PrivacySettings) SetOnlineStatus(s PrivacySetting)     { p.onlineStatus = s }
func (p *PrivacySettings) SetFirstDialogsInit(s PrivacySetting) { p.firstDialogsInit = s }
func (p *PrivacySettings) SetGroupChatInvites(s PrivacySetting) { p.groupChatInvites = s }
