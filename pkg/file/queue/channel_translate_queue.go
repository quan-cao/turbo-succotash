package queue

import "time"

type ChannelTranslateQueue struct {
	c chan *TranslateTask
}

func NewChannelTranslateQueue(c chan *TranslateTask) *ChannelTranslateQueue {
	return &ChannelTranslateQueue{c}
}

func (q *ChannelTranslateQueue) Add(t *TranslateTask) error {
	q.c <- t
	return nil
}

func (q *ChannelTranslateQueue) Take() (*TranslateTask, string) {
	select {
	case t := <-q.c:
		return t, ""
	case <-time.After(100 * time.Millisecond):
		return nil, ""
	}
}

func (q *ChannelTranslateQueue) Delete(key string) error {
	return nil
}

// Ensure implementation
var _ TranslateQueue = (*ChannelTranslateQueue)(nil)
