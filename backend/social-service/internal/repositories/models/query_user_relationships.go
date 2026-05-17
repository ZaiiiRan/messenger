package models

import userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"

type QueryUserRelationshipsDal struct {
	FirstUserId   *string  `json:"first_user_id"`
	SecondUserIds []string `json:"second_user_ids"`

	Statuses []int16 `json:"statuses"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`

	OrderByUpdatedAtDesc bool `json:"order_by_updated_at_desc"`
}

func NewQueryUserRelationshipsDal(
	firstUserId *string,
	secondUserIds []string,
	statuses []userrelationship.UserRelationshipStatus,
	page int, pageSize int,
	orderByUpdatedAtDesc bool,
) *QueryUserRelationshipsDal {
	convertedStatuses := make([]int16, len(statuses))
	for i, status := range statuses {
		convertedStatuses[i] = int16(status)
	}

	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}

	return &QueryUserRelationshipsDal{
		FirstUserId:          firstUserId,
		SecondUserIds:        secondUserIds,
		Statuses:             convertedStatuses,
		Limit:                pageSize,
		Offset:               (page - 1) * pageSize,
		OrderByUpdatedAtDesc: orderByUpdatedAtDesc,
	}
}
