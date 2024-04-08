package entity

import "time"

type OriginalFileMetadata struct {
	Id             int
	SHA256         string
	Filename       string
	FileType       string
	FileSize       int
	SourceLanguage string
	TokenCount     int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      string
}
