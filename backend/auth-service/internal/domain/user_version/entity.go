package userversion

import "time"

type UserVersion struct {
	id        int64
	userId    string
	version   int
	createdAt time.Time
	updatedAt time.Time
}

func New(userId string) *UserVersion {
	now := time.Now()

	return &UserVersion{
		userId:    userId,
		version:   1,
		createdAt: now,
		updatedAt: now,
	}
}

func FromStorage(
	id int64,
	userId string,
	version int,
	createdAt, updatedAt time.Time,
) *UserVersion {
	return &UserVersion{
		id:        id,
		userId:    userId,
		version:   version,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (uv *UserVersion) Id() int64            { return uv.id }
func (uv *UserVersion) UserId() string       { return uv.userId }
func (uv *UserVersion) Version() int         { return uv.version }
func (uv *UserVersion) CreatedAt() time.Time { return uv.createdAt }
func (uv *UserVersion) UpdatedAt() time.Time { return uv.updatedAt }

func (uv *UserVersion) SetId(id int64) {
	if uv.id == 0 {
		uv.id = id
	}
}

func (uv *UserVersion) IncrementVersion() {
	uv.version++
	uv.updatedAt = time.Now()
}
