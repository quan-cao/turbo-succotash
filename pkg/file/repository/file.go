package repository

// FileRepository operates against a filesystem as it
// needs to store and retrieve binary files.
type FileRepository interface {
	Persist(b []byte, filepath string) error
	Get(filepath string) ([]byte, error)
	Delete(filepath string) error
	DeleteMany(filepaths []string) error
	GetUrl(filepath string) (string, error)
}
