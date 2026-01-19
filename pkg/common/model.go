package common

import "time"

type UploadFile struct {
	ID        int64     `database:"id" json:"id"`
	Name      string    `database:"name" json:"name"`
	S3Path    string    `database:"s3_path" json:"s3_path"`
	CreatedAt time.Time `database:"created_at" json:"created_at"`
	UpdatedAt time.Time `database:"updated_at" json:"updated_at"`
}
