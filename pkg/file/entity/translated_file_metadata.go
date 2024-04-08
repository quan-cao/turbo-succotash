package entity

import "time"

type TranslatedFileMetadata struct {
	Id             int
	OriginalFileId int
	Filename       string
	TargetLanguage string
	Cost           float64
	TimeTaken      int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      string
}
