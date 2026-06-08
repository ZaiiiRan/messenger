package models

import (
	"time"

	userphoto "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user_photo"
)

type V1UserPhotoDal struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	FileId    string    `db:"file_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func V1UserPhotoDalFromDomain(p *userphoto.UserPhoto) V1UserPhotoDal {
	return V1UserPhotoDal{
		Id:        p.GetID(),
		UserId:    p.GetUserID(),
		FileId:    p.GetFileID(),
		CreatedAt: p.GetCreatedAt(),
		UpdatedAt: p.GetUpdatedAt(),
	}
}

func (d V1UserPhotoDal) ToDomain() *userphoto.UserPhoto {
	return userphoto.FromStorage(d.Id, d.UserId, d.FileId, d.CreatedAt, d.UpdatedAt)
}
