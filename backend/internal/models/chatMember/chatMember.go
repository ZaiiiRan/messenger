package chatMember

import (
	"backend/internal/models/shortUser"
)

type ChatMember struct {
	User      *shortUser.ShortUser `json:"user"`
	Role      int                  `json:"role"`
	RemovedBy *uint64
	AddedBy   uint64
	ChatID    uint64 `json:"chat_id"`
}

func (member *ChatMember) Removed() bool {
	return member.RemovedBy != nil
}