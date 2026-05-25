package models

type FileDeletionTask struct {
	Id      string   `json:"id"`
	FileIds []string `json:"file_ids"`
}
