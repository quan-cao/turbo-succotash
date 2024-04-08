package queue

import (
	"doc-translate-go/mocks"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/mock/gomock"
)

func newMock(t *testing.T, queueUrl string, groupId string) (*mocks.MockSQSAPI, *SqsTranslateQueue) {
	ctrl := gomock.NewController(t)
	mockedClient := mocks.NewMockSQSAPI(ctrl)
	queue := NewSqsTranslateQueue(mockedClient, queueUrl, groupId)
	return mockedClient, queue
}

func TestSqsTranslateQueue__Add(t *testing.T) {
	mockedClient, queue := newMock(t, "url", "id")

	task := &TranslateTask{
		Isid:           "1",
		Filename:       "filename",
		SourceLang:     "sourcelang",
		TargetLang:     "targetlang",
		OriginalFileId: 1,
	}

	taskJson, _ := json.Marshal(task)
	message := &sqs.SendMessageInput{
		QueueUrl:       aws.String(queue.queueUrl),
		MessageBody:    aws.String(string(taskJson)),
		MessageGroupId: aws.String(queue.groupId),
	}

	mockedClient.EXPECT().SendMessage(gomock.Eq(message)).Times(1)

	err := queue.Add(task)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSqsTranslateQueue__Add_SendMessageErr(t *testing.T) {
	mockedClient, queue := newMock(t, "url", "id")

	task := &TranslateTask{
		Isid:           "1",
		Filename:       "filename",
		SourceLang:     "sourcelang",
		TargetLang:     "targetlang",
		OriginalFileId: 1,
	}

	taskJson, _ := json.Marshal(task)
	message := &sqs.SendMessageInput{
		QueueUrl:       aws.String(queue.queueUrl),
		MessageBody:    aws.String(string(taskJson)),
		MessageGroupId: aws.String(queue.groupId),
	}

	e := errors.New("message error")
	mockedClient.EXPECT().SendMessage(gomock.Eq(message)).Times(1).Return(nil, e)

	err := queue.Add(task)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
}

func TestSqsTranslateQueue__Take(t *testing.T) {
	mockedClient, queue := newMock(t, "url", "id")

	task := &TranslateTask{
		Isid:           "1",
		Filename:       "filename",
		SourceLang:     "sourcelang",
		TargetLang:     "targetlang",
		OriginalFileId: 1,
	}

	taskJson, _ := json.Marshal(task)
	wantSendMessageParam := &sqs.SendMessageInput{
		QueueUrl:       aws.String(queue.queueUrl),
		MessageBody:    aws.String(string(taskJson)),
		MessageGroupId: aws.String(queue.groupId),
	}

	mockedClient.EXPECT().SendMessage(gomock.Eq(wantSendMessageParam)).Times(1)
	err := queue.Add(task)
	if err != nil {
		t.Fatal(err)
	}

	receiveKey := "1"
	output := &sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			{Body: aws.String(string(taskJson)), ReceiptHandle: aws.String(receiveKey)},
		},
	}

	wantReceiveMessageParams := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(1),
		QueueUrl:            aws.String("url"),
		WaitTimeSeconds:     aws.Int64(2),
	}
	mockedClient.EXPECT().ReceiveMessage(gomock.Eq(wantReceiveMessageParams)).Times(1).Return(output, nil)

	got, key := queue.Take()
	if !reflect.DeepEqual(task, got) {
		t.Fatalf("expected %v, got %v", task, got)
	}
	if key != receiveKey {
		t.Fatalf("expected %v, got %v", receiveKey, key)
	}
}

func TestSqsTranslateQueue__Take_RetrieveMessageErr(t *testing.T) {
	mockedClient, queue := newMock(t, "url", "id")

	task := &TranslateTask{
		Isid:           "1",
		Filename:       "filename",
		SourceLang:     "sourcelang",
		TargetLang:     "targetlang",
		OriginalFileId: 1,
	}

	taskJson, _ := json.Marshal(task)
	wantSendMessageParam := &sqs.SendMessageInput{
		QueueUrl:       aws.String(queue.queueUrl),
		MessageBody:    aws.String(string(taskJson)),
		MessageGroupId: aws.String(queue.groupId),
	}

	mockedClient.EXPECT().SendMessage(gomock.Eq(wantSendMessageParam)).Times(1)
	err := queue.Add(task)
	if err != nil {
		t.Fatal(err)
	}

	wantReceiveMessageParams := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(1),
		QueueUrl:            aws.String("url"),
		WaitTimeSeconds:     aws.Int64(2),
	}

	e := errors.New("error")

	mockedClient.
		EXPECT().
		ReceiveMessage(gomock.Eq(wantReceiveMessageParams)).
		Times(1).
		Return(nil, e)

	got, key := queue.Take()
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
	if key != "" {
		t.Fatalf("expected empty key, got %v", key)
	}
}

func TestSqsTranslateQueue__Delete(t *testing.T) {
	mockedClient, queue := newMock(t, "url", "id")

	wantDeleteMessageParam := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queue.queueUrl),
		ReceiptHandle: aws.String("1"),
	}

	mockedClient.EXPECT().DeleteMessage(gomock.Eq(wantDeleteMessageParam)).Times(1)

	err := queue.Delete("1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSqsTranslateQueue__Delete_DeleteMessageErr(t *testing.T) {
	mockedClient, queue := newMock(t, "url", "id")

	wantDeleteMessageParam := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queue.queueUrl),
		ReceiptHandle: aws.String("1"),
	}

	e := errors.New("error")
	mockedClient.EXPECT().DeleteMessage(gomock.Eq(wantDeleteMessageParam)).Times(1).Return(nil, e)

	err := queue.Delete("1")
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
}
