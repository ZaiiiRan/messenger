package models

type QueryExpiredCodesDal struct {
	PageSize int `json:"page_size"`
}

func NewQueryExpiredCodesDal(pageSize int) *QueryExpiredCodesDal {
	return &QueryExpiredCodesDal{
		PageSize: pageSize,
	}
}
