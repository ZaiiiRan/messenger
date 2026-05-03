package models

type QueryTokenDal struct {
	Id      *int64  `json:"id"`
	UserId  *string `json:"user_id"`
	Token   *string `json:"token"`
	Version *int    `json:"version"`
}

func NewQueryTokenDal(id *int64, userId *string, token *string, version *int) *QueryTokenDal {
	return &QueryTokenDal{
		Id:      id,
		UserId:  userId,
		Token:   token,
		Version: version,
	}
}
