package queue

import (
	"reflect"
	"testing"
)

func newQueue() *ChannelTranslateQueue {
	c := make(chan *TranslateTask, 1<<32)
	return NewChannelTranslateQueue(c)
}

func TestChannelTranslateQueue(t *testing.T) {
	q := newQueue()

	want := &TranslateTask{}

	err := q.Add(want)
	if err != nil {
		t.Fatal(err)
	}

	got, _ := q.Take()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expect %v, got %v", want, got)
	}

	got, _ = q.Take()
	if got != nil {
		t.Fatalf("expect nil, got %v", got)
	}
}
