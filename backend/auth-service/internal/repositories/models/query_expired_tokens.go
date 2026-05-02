package models

type QueryExpiredTokensDal struct {
	PageSize int `json:"page_size"`
}

func NewQueryExpiredTokensDal(pageSize int) *QueryExpiredTokensDal {
	return &QueryExpiredTokensDal{
		PageSize: pageSize,
	}
}
