package translator

type Translator interface {
	Translate(b []byte, sourceLang string, targetLang string) ([]byte, error)
}
