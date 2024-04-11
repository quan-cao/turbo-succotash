package tracker

import "time"

type FileTrackerInput struct {
	Status     string    `json:"status"`
	SourceLang string    `json:"source_lang"`
	TargetLang string    `json:"target_lang"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FileStatus struct {
	Key        string    `json:"key"`
	Status     string    `json:"status"`
	SourceLang string    `json:"source_lang"`
	TargetLang string    `json:"target_lang"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FileTracker interface {
	Get(key string) (*FileStatus, error)
	Create(status *FileStatus) error
	Delete(key string) error
	List(pat string) ([]*FileStatus, error)
	Clear() error
}
