package chat

import (
	"backend/internal/dtos/userDTO"
)

type Chat struct {
	ID          uint64       `json:"id"`
	Name        string       `json:"name"`
	IsGroupChat bool         `json:"is_group_chat"`
	IsDeleted   bool         `json:"is_deleted"`
	Members     []ChatMember `json:"members"`
}

type ChatMember struct {
	User userDTO.UserDTO `json:"user"`
	Role string          `json:"role"`
}