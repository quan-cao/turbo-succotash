package queue

type TranslateTask struct {
	Isid           string `json:"isid"`
	Filename       string `json:"filename"`
	SourceLang     string `json:"source_lang"`
	TargetLang     string `json:"target_lang"`
	OriginalFileId int    `json:"original_file_id"`
}

type TranslateQueue interface {
	Add(t *TranslateTask) error
	Take() (*TranslateTask, string)
	Delete(key string) error
}
