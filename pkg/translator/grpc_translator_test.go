package translator

import (
	"doc-translate-go/mocks"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGrpcTranslator(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := mocks.NewMockDocumentProcessorClient(ctrl)

	translatr := NewGrpcTranslator(c)

	c.EXPECT().ProcessDocument(gomock.Any(), gomock.Any()).Times(1)

	_, err := translatr.Translate([]byte{}, "sourceLang", "targetLang")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGrpcTranslator__Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := mocks.NewMockDocumentProcessorClient(ctrl)

	translatr := NewGrpcTranslator(c)

	e := errors.New("process error")
	c.EXPECT().ProcessDocument(gomock.Any(), gomock.Any()).Return(nil, e).Times(1)

	b, err := translatr.Translate([]byte{}, "sourceLang", "targetLang")
	if err == nil {
		t.Fatal("expect error, got nil")
	}

	if b != nil {
		t.Fatalf("expect nil, got %v", b)
	}
}
