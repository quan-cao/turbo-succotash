package usecase

import (
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

	uc.fileTracker.Add(isid, filename, sourceLang, targetLang, "in progress")

	metadatas, err := uc.originalFileMetaUC.ListByFilenameIsid(filename, isid)
	if err != nil {
		return errors.New("failed to check for duplicated files")
	}

	if len(metadatas) > 0 {
		uc.fileTracker.Add(isid, filename, sourceLang, targetLang, "fail:duplicate")
		return errors.New("failed to check for duplicated files")
	}

	err = uc.fileUC.Persist(b, fmt.Sprintf("%s/%s", isid, filename))
	if err != nil {
		uc.fileTracker.Add(isid, filename, sourceLang, targetLang, "fail:persist")
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

func (uc *TranslateUseCase) ListenAndExecute() {
	for range time.Tick(500 * time.Millisecond) {
		t, key := uc.translateQueue.Take()

		b, err := uc.fileUC.Get(fmt.Sprintf("%s/%s", t.Isid, t.Filename))
		if err != nil {
			uc.fileTracker.Add(t.Isid, t.Filename, t.SourceLang, t.TargetLang, "fail:read")
			continue
		}

		translated_b, err := uc.Translate(b, t.SourceLang, t.TargetLang)
		if err != nil {
			uc.fileTracker.Add(t.Isid, t.Filename, t.SourceLang, t.TargetLang, "fail:translate")
			continue
		}

		translatedFilename := fmt.Sprintf("translated-%s-to-%s-%s", t.SourceLang, t.TargetLang, t.Filename)
		err = uc.fileUC.Persist(translated_b, fmt.Sprintf("%s/%s", t.Isid, translatedFilename))
		if err != nil {
			uc.fileTracker.Add(t.Isid, translatedFilename, t.SourceLang, t.TargetLang, "fail:persist")
			continue
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

		uc.fileTracker.Clear(t.Isid, t.Filename)

		uc.translateQueue.Delete(key)
	}
}
