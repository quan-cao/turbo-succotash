package translator

// GrpcTranslator makes gRPC call to another service to translate documents
type EchoTranslator struct{}

func NewEchoTranslator() *EchoTranslator {
	return &EchoTranslator{}
}

func (t *EchoTranslator) Translate(b []byte, sourceLang string, targetLang string) ([]byte, error) {
	return b, nil
}

// Ensure implementation
var _ Translator = (*EchoTranslator)(nil)
