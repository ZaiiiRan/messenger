package chatMemberDTO

import (
	"backend/internal/models/chatMember"
	"backend/internal/models/shortUser"
)

type ChatMemberDTO struct {
	User   *shortUser.ShortUser `json:"user"`
	Role   string               `json:"role"`
	ChatID uint64               `json:"chat_id"`
}

func CreateChatMemberDTO(member *chatMember.ChatMember) *ChatMemberDTO {
	role := chatMember.GetRoleString(member.Role)
	return &ChatMemberDTO{
		User:   member.User,
		Role:   role,
		ChatID: member.ChatID,
	}
}
