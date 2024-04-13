package usecase

import (
	"context"
	"doc-translate-go/pkg/file/entity"
	"doc-translate-go/pkg/file/queue"
	"doc-translate-go/pkg/tracker"
	"doc-translate-go/pkg/translator"
	"errors"
	"fmt"
	"path/filepath"
	"time"
)

type TranslateUseCase struct {
	translator           translator.Translator
	originalFileMetaUC   *OriginalFileMetadataUseCase
	translatedFileMetaUC *TranslatedFileMetadataUseCase
	fileUC               *FileUseCase
	fileTracker          tracker.FileTracker
	translateQueue       queue.TranslateQueue
}

func NewTranslateUseCase(
	translator translator.Translator,
	originalFileMetadataUseCase *OriginalFileMetadataUseCase,
	translateFileMetadataUseCase *TranslatedFileMetadataUseCase,
	fileUseCase *FileUseCase,
	fileTracker tracker.FileTracker,
	translateQueue queue.TranslateQueue,
) *TranslateUseCase {
	return &TranslateUseCase{
		translator,
		originalFileMetadataUseCase,
		translateFileMetadataUseCase,
		fileUseCase,
		fileTracker,
		translateQueue,
	}
}

func (uc *TranslateUseCase) Translate(b []byte, sourceLang string, targetLang string) ([]byte, error) {
	return uc.translator.Translate(b, sourceLang, targetLang)
}

// TranslateAsync stores file in filesystem and sends a message to a queue
func (uc *TranslateUseCase) TranslateAsync(
	b []byte,
	filename string,
	filesize int,
	isid string,
	sourceLang string,
	targetLang string,
) error {
	fileExt := filepath.Ext(filename)
	if fileExt != "docx" {
		return errors.New("invalid extension")
	}

	fileTrackerKey := fmt.Sprintf("%s_%s", isid, filename)

	uc.fileTracker.Create(&tracker.FileStatus{
		Key:        fileTrackerKey,
		Status:     "in progress",
		SourceLang: sourceLang,
		TargetLang: targetLang,
	})

	metadatas, err := uc.originalFileMetaUC.ListByFilenameIsid(filename, isid)
	if err != nil {
		return errors.New("failed to check for duplicated files")
	}

	if len(metadatas) > 0 {
		uc.fileTracker.Create(&tracker.FileStatus{
			Key:        fileTrackerKey,
			Status:     "fail:duplicate",
			SourceLang: sourceLang,
			TargetLang: targetLang,
		})
		return errors.New("failed to check for duplicated files")
	}

	err = uc.fileUC.Persist(b, fmt.Sprintf("%s/%s", isid, filename))
	if err != nil {
		uc.fileTracker.Create(&tracker.FileStatus{
			Key:        fileTrackerKey,
			Status:     "fail:persist",
			SourceLang: sourceLang,
			TargetLang: targetLang,
		})
		return errors.New("failed to persist file")
	}

	fileType := "docx"
	now := time.Now()
	id, err := uc.originalFileMetaUC.Persist(&entity.OriginalFileMetadata{
		SHA256:         "",
		Filename:       filename,
		FileType:       fileType,
		FileSize:       filesize,
		SourceLanguage: sourceLang,
		TokenCount:     0,
		CreatedAt:      now,
		UpdatedAt:      now,
		CreatedBy:      isid,
	})
	if err != nil {
		return errors.New("failed to store original file metadata")
	}

	err = uc.translateQueue.Add(&queue.TranslateTask{
		Isid:           isid,
		Filename:       filename,
		SourceLang:     sourceLang,
		TargetLang:     targetLang,
		OriginalFileId: id,
	})
	if err != nil {
		return err
	}

	return nil
}

// ExecuteQueue takes one message out of the queue and performs translating.
func (uc *TranslateUseCase) ExecuteQueue() error {
	t, key := uc.translateQueue.Take()
	if t == nil {
		return nil
	}

	fileTrackerKey := fmt.Sprintf("%s_%s", t.Isid, t.Filename)

	b, err := uc.fileUC.Get(fmt.Sprintf("%s/%s", t.Isid, t.Filename))
	if err != nil {
		uc.fileTracker.Create(&tracker.FileStatus{
			Key:        fileTrackerKey,
			Status:     "fail:read",
			SourceLang: t.SourceLang,
			TargetLang: t.TargetLang,
		})
		return err
	}

	translated_b, err := uc.Translate(b, t.SourceLang, t.TargetLang)
	if err != nil {
		uc.fileTracker.Create(&tracker.FileStatus{
			Key:        fileTrackerKey,
			Status:     "fail:translate",
			SourceLang: t.SourceLang,
			TargetLang: t.TargetLang,
		})
		return err
	}

	translatedFilename := fmt.Sprintf("translated-%s-to-%s-%s", t.SourceLang, t.TargetLang, t.Filename)
	err = uc.fileUC.Persist(translated_b, fmt.Sprintf("%s/%s", t.Isid, translatedFilename))
	if err != nil {
		uc.fileTracker.Create(&tracker.FileStatus{
			Key:        fileTrackerKey,
			Status:     "fail:persist",
			SourceLang: t.SourceLang,
			TargetLang: t.TargetLang,
		})
		return err
	}

	now := time.Now()
	uc.translatedFileMetaUC.Persist(&entity.TranslatedFileMetadata{
		OriginalFileId: t.OriginalFileId,
		Filename:       translatedFilename,
		TargetLanguage: t.TargetLang,
		Cost:           0,
		CreatedAt:      now,
		UpdatedAt:      now,
		CreatedBy:      t.Isid,
	})

	uc.fileTracker.Delete(fileTrackerKey)

	uc.translateQueue.Delete(key)

	return nil
}

// ListenAndExecute polls from translate queue and performs translating one message by another.
func (uc *TranslateUseCase) ListenAndExecute(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			uc.ExecuteQueue()
		}
	}
}
