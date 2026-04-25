package models

type QueryPasswordDal struct {
	Id     *int64  `json:"id"`
	UserId *string `json:"user_id"`
}

func NewQueryPasswordDal(id *int64, userId *string) *QueryPasswordDal {
	return &QueryPasswordDal{
		Id:     id,
		UserId: userId,
	}
}
