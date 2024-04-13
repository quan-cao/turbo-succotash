package translator

import (
	"reflect"
	"testing"
)

func TestEchoTranslator(t *testing.T) {
	translatr := NewEchoTranslator()

	want := []byte{}

	got, err := translatr.Translate([]byte{}, "sourceLang", "targetLang")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
