package tracker

import "time"

type FileProgress struct {
	Status     *string    `json:"status"`
	SourceLang *string    `json:"source_lang"`
	TargetLang *string    `json:"target_lang"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type FileTracker interface {
	Add(isid string, filename string, sourceLang string, targetLang string, status string) error
	Clear(isid string, filename string) error
}
