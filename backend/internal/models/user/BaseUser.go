package user

type BaseUser struct {
	ID          uint64 `json:"user_id"`
	Username    string `json:"username"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	IsDeleted   bool   `json:"is_deleted"`
	IsBanned    bool   `json:"is_banned"`
	IsActivated bool   `json:"is_activated"`
}

// New Base User
func NewBaseUser(username, firstname, lastname string) BaseUser {
	return BaseUser{
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
	}
}
