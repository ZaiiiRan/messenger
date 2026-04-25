package models

type QueryCodeDal struct {
	Id     *int64  `json:"id"`
	UserId *string `json:"user_id"`
}

func NewQueryCodeDal(id *int64, userId *string) *QueryCodeDal {
	return &QueryCodeDal{
		Id:     id,
		UserId: userId,
	}
}
