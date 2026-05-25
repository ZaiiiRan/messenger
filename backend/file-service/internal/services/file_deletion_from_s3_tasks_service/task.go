package filedeletionfroms3tasksservice

type FileDeletionFromS3Task struct {
	Bucket int16    `json:"bucket"`
	S3Keys []string `json:"s3_keys"`
}
