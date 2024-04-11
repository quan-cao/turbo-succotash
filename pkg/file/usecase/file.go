package usecase

import (
	"doc-translate-go/pkg/file/repository"
	"sync"
)

type FileUseCase struct {
	repo repository.FileRepository
}

func NewFileUseCase(repo repository.FileRepository) *FileUseCase {
	return &FileUseCase{repo}
}

func (uc *FileUseCase) Get(filepath string) ([]byte, error) {
	return uc.repo.Get(filepath)
}

func (uc *FileUseCase) GetMany(filepaths []string) (map[string][]byte, map[string]error) {
	wg := sync.WaitGroup{}

	data := make(map[string][]byte)
	errors := make(map[string]error)

	for _, p := range filepaths {
		wg.Add(1)

		go func() {
			defer wg.Done()

			dat, err := uc.Get(p)
			data[p] = dat
			errors[p] = err
		}()
	}

	wg.Wait()

	return data, errors
}

func (uc *FileUseCase) Persist(b []byte, filepath string) error {
	return uc.repo.Persist(b, filepath)
}

func (uc *FileUseCase) Delete(filepath string) error {
	return uc.repo.Delete(filepath)
}

func (uc *FileUseCase) DeleteMany(filepaths []string) error {
	return uc.repo.DeleteMany(filepaths)
}

func (uc *FileUseCase) GetUrls(filepaths []string) (map[string]string, map[string]error) {
	urls := make(map[string]string)
	errors := make(map[string]error)

	for _, filepath := range filepaths {
		url, err := uc.repo.GetUrl(filepath)
		urls[filepath] = url
		errors[filepath] = err
	}

	return urls, errors
}
