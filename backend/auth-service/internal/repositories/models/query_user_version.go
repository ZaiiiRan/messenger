package models

type QueryUserVersionDal struct {
	Id     *int64  `json:"id"`
	UserId *string `json:"user_id"`
}

func NewQueryUserVersionDal(id *int64, userId *string) *QueryUserVersionDal {
	return &QueryUserVersionDal{
		Id:     id,
		UserId: userId,
	}
}
