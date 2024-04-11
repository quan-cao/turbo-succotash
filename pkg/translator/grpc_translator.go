package translator

import (
	"context"
	"time"

	documentprotov1 "doc-translate-go/proto/gen/go/proto/documentprocessor/v1"
)

// GrpcTranslator makes gRPC call to another service to translate documents
type GrpcTranslator struct {
	client documentprotov1.DocumentProcessorClient
}

func NewGrpcTranslator(client documentprotov1.DocumentProcessorClient) *GrpcTranslator {
	return &GrpcTranslator{client}
}

func (t *GrpcTranslator) Translate(b []byte, sourceLang string, targetLang string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Second)
	defer cancel()

	resp, err := t.client.ProcessDocument(ctx, &documentprotov1.DocumentRequest{
		Document:   b,
		SourceLang: sourceLang,
		TargetLang: targetLang,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetDocument(), nil
}

// Ensure implementation
var _ Translator = (*GrpcTranslator)(nil)
