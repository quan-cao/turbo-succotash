package queue

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type SqsTranslateQueue struct {
	client   sqsiface.SQSAPI
	queueUrl string
	groupId  string
}

func NewSqsTranslateQueue(client sqsiface.SQSAPI, queueUrl string, groupId string) *SqsTranslateQueue {
	return &SqsTranslateQueue{client, queueUrl, groupId}
}

func (q *SqsTranslateQueue) Add(t *TranslateTask) error {
	taskJson, err := json.Marshal(t)
	if err != nil {
		return err
	}

	_, err = q.client.SendMessage(&sqs.SendMessageInput{
		QueueUrl:       aws.String(q.queueUrl),
		MessageBody:    aws.String(string(taskJson)),
		MessageGroupId: aws.String(q.groupId),
	})
	if err != nil {
		return err
	}

	return nil
}

func (q *SqsTranslateQueue) Take() (*TranslateTask, string) {
	result, err := q.client.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.queueUrl),
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(2),
	})
	if err != nil {
		return nil, ""
	}

	var task *TranslateTask
	message := result.Messages[0]
	err = json.Unmarshal([]byte(*message.Body), &task)
	if err != nil {
		return nil, ""
	}

	return task, *message.ReceiptHandle
}

func (q *SqsTranslateQueue) Delete(key string) error {
	_, err := q.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.queueUrl),
		ReceiptHandle: aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}

// Ensure implementation
var _ TranslateQueue = (*SqsTranslateQueue)(nil)
