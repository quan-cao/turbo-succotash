package translator

import (
	"context"
	"time"

	documentproto "doc-translate-go/gen/go/proto/documentprocessor"
)

// GrpcTranslator makes gRPC call to another service to translate documents
type GrpcTranslator struct {
	client documentproto.DocumentProcessorClient
}

func NewGrpcTranslator(client documentproto.DocumentProcessorClient) *GrpcTranslator {
	return &GrpcTranslator{client}
}

func (t *GrpcTranslator) Translate(b []byte, sourceLang string, targetLang string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Second)
	defer cancel()

	resp, err := t.client.ProcessDocument(
		ctx,
		&documentproto.DocumentRequest{
			Document:   b,
			SourceLang: sourceLang,
			TargetLang: targetLang,
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.GetDocument(), nil
}

// Ensure implementation
var _ Translator = (*GrpcTranslator)(nil)
