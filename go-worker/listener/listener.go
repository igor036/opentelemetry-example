package listener

import (
	"go-worker/internal/settings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	svc            *sqs.SQS
	sqsMessagePoll chan *sqs.Message
)

func DispatchMessage(timeout int64, queueName string, HandlerSQSMessage func(message *sqs.Message)) {

	go ReciverSQSMessages(timeout, queueName)

	for message := range sqsMessagePoll {
		HandlerSQSMessage(message)
		DeleteSQSMessage(queueName, message)
	}
}

func ReciverSQSMessages(timeout int64, queueName string) {

	queueURL := sqsQueueURL(&queueName).QueueUrl
	visibilityTimeout := sqsVisibilityTimout(timeout)

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(settings.Env.SQS.MaxNumberOfMessages),
		VisibilityTimeout:   aws.Int64(visibilityTimeout),
	})

	if err != nil {
		panic(err)
	}

	for _, message := range result.Messages {
		sqsMessagePoll <- message
	}
}

func DeleteSQSMessage(queueName string, message *sqs.Message) {

	queueURL := sqsQueueURL(&queueName).QueueUrl

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: message.ReceiptHandle,
	})

	if err != nil {
		panic(err)
	}
}

func awsSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

func sqsQueueURL(queueName *string) *sqs.GetQueueUrlOutput {

	queueURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queueName,
	})

	if err != nil {
		panic(err)
	}

	return queueURL
}

func sqsVisibilityTimout(timeout int64) int64 {
	visibilityTimeout := time.Duration(timeout) * time.Second
	return int64(visibilityTimeout.Seconds())
}

func init() {
	svc = sqs.New(awsSession())
	sqsMessagePoll = make(chan *sqs.Message, settings.Env.SQS.MaxNumberOfMessages)
}
