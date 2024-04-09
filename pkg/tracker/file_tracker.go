package tracker

import "time"

type FileStatus struct {
	Status     string    `json:"status"`
	SourceLang string    `json:"source_lang"`
	TargetLang string    `json:"target_lang"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FileTracker interface {
	Get(key string) (*FileStatus, error)
	Create(key string, fileProgress *FileStatus) error
	Delete(key string) error
	List(pat string) ([]*FileStatus, error)
	Clear() error
}
