package models

import "time"

type UserRelationshipChangeTask struct {
	Id                     string    `json:"id"`
	User1Id                string    `json:"user1_id"`
	User2Id                string    `json:"user2_id"`
	UserRelationshipStatus string    `json:"user_relationship_status"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
