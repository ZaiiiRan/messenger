package userphoto

import "time"

type UserPhoto struct {
	id        string
	userID    string
	fileID    string
	createdAt time.Time
	updatedAt time.Time
}

func New(userID, fileID string) *UserPhoto {
	now := time.Now()
	return &UserPhoto{
		userID:    userID,
		fileID:    fileID,
		createdAt: now,
		updatedAt: now,
	}
}

func FromStorage(id, userID, fileID string, createdAt, updatedAt time.Time) *UserPhoto {
	return &UserPhoto{
		id:        id,
		userID:    userID,
		fileID:    fileID,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (p *UserPhoto) GetID() string        { return p.id }
func (p *UserPhoto) GetUserID() string    { return p.userID }
func (p *UserPhoto) GetFileID() string    { return p.fileID }
func (p *UserPhoto) GetCreatedAt() time.Time { return p.createdAt }
func (p *UserPhoto) GetUpdatedAt() time.Time { return p.updatedAt }

func (p *UserPhoto) SetID(id string) {
	if p.id == "" {
		p.id = id
	}
}
