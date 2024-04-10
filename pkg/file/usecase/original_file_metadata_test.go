package usecase

import (
	"doc-translate-go/mocks"
	"doc-translate-go/pkg/file/entity"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func newMock(t *testing.T) (*gomock.Controller, *mocks.MockOriginalFileMetadataRepository) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockOriginalFileMetadataRepository(ctrl)
	return ctrl, repo
}

func TestOriginalFileMetadataUseCase__Persist(t *testing.T) {
	_, repo := newMock(t)
	uc := NewOriginalFileMetadataUseCase(repo)

	ent := &entity.OriginalFileMetadata{}
	want := 1

	repo.EXPECT().Create(gomock.Eq(ent)).Times(1).Return(want, nil)

	id, err := uc.Persist(ent)
	if err != nil {
		t.Fatal(err)
	}

	if id != want {
		t.Fatalf("expected %v, got %v", want, id)
	}
}

func TestOriginalFileMetadataUseCase__Persist_Err(t *testing.T) {
	_, repo := newMock(t)
	uc := NewOriginalFileMetadataUseCase(repo)

	ent := &entity.OriginalFileMetadata{}

	e := errors.New("create error")
	repo.EXPECT().Create(gomock.Eq(ent)).Times(1).Return(0, e)

	id, err := uc.Persist(ent)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if id != 0 {
		t.Fatalf("expected %v, got %v", 0, id)
	}
}

func TestOriginalFileMetadataUseCase__ListByFilenameIsid(t *testing.T) {
	_, repo := newMock(t)
	uc := NewOriginalFileMetadataUseCase(repo)

	want := []*entity.OriginalFileMetadata{{Id: 1}, {Id: 2}}

	repo.EXPECT().ListByFilenameIsid(gomock.Eq("filename"), gomock.Eq("isid")).Times(1).Return(want, nil)

	got, err := uc.ListByFilenameIsid("filename", "isid")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestOriginalFileMetadataUseCase__ListByFilenameIsid_Err(t *testing.T) {
	_, repo := newMock(t)
	uc := NewOriginalFileMetadataUseCase(repo)

	e := errors.New("list error")
	repo.EXPECT().ListByFilenameIsid(gomock.Eq("filename"), gomock.Eq("isid")).Times(1).Return(nil, e)

	got, err := uc.ListByFilenameIsid("filename", "isid")
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}
