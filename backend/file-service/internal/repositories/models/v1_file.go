package models

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
)

type V1FileDal struct {
	Id        string    `json:"id" db:"id"`
	OwnerId   string    `json:"owner_id" db:"owner_id"`
	Status    int16     `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func V1FileFromDomain(file *file.File) V1FileDal {
	if file == nil {
		return V1FileDal{}
	}
	return V1FileDal{
		Id:        file.GetID(),
		OwnerId:   file.GetOwnerID(),
		Status:    int16(file.GetStatus()),
		CreatedAt: file.GetCreatedAt(),
		UpdatedAt: file.GetUpdatedAt(),
	}
}

func (f V1FileDal) IsNull() bool { return false }
func (f V1FileDal) Index(i int) any {
	switch i {
	case 0:
		return f.Id
	case 1:
		return f.OwnerId
	case 2:
		return f.Status
	case 3:
		return f.CreatedAt
	case 4:
		return f.UpdatedAt
	default:
		return nil
	}
}

func (f V1FileDal) ToDomain() *file.File {
	return file.FromStorage(
		f.Id,
		f.OwnerId,
		f.Status,
		f.CreatedAt,
		f.UpdatedAt,
	)
}
