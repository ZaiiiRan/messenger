package chatMemberDTO

import (
	"backend/internal/models/chat/chatMember"
	"backend/internal/models/shortUser"
)

type ChatMemberDTO struct {
	User      *shortUser.ShortUser `json:"user"`
	Role      string               `json:"role"`
	ChatID    uint64               `json:"chat_id"`
	IsRemoved bool                 `json:"is_removed"`
	IsLeft    bool                 `json:"is_left"`
	AddedBy   uint64               `json:"added_by"`
}

func CreateChatMemberDTO(member *chatMember.ChatMember) *ChatMemberDTO {
	role := chatMember.GetRoleString(member.Role)
	return &ChatMemberDTO{
		User:      member.User,
		Role:      role,
		ChatID:    member.ChatID,
		IsRemoved: member.IsRemoved(),
		IsLeft:    member.IsLeft(),
		AddedBy:   member.AddedBy,
	}
}

// Converting chat member array to chat member dto array
func GetChatMembersDTOs(members []chatMember.ChatMember) []*ChatMemberDTO {
	chatMembersDTOs := make([]*ChatMemberDTO, len(members))
	for index, member := range members {
		chatMembersDTOs[index] = CreateChatMemberDTO(&member)
	}
	return chatMembersDTOs
}