package models

type QueryTokensDal struct {
	UserId       string
	ExcludeToken string
	Version      int
	Limit        int
	Offset       int
}

func NewQueryTokensDal(userId, excludeToken string, version int, page, pageSize int) *QueryTokensDal {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return &QueryTokensDal{
		UserId:       userId,
		ExcludeToken: excludeToken,
		Version:      version,
		Limit:        pageSize,
		Offset:       (page - 1) * pageSize,
	}
}
