package harvest

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
)

var CommentQueueUrl = "https://sqs.us-east-2.amazonaws.com/583491089160/manessingercomcomments.fifo"
var WaitTimeSeconds int64 = 2
var MaxNumberOfMessages int64 = 10

var ReceiveMessageInput = sqs.ReceiveMessageInput{
	QueueUrl:            &CommentQueueUrl,
	WaitTimeSeconds:     &WaitTimeSeconds,
	MaxNumberOfMessages: &MaxNumberOfMessages,
}

var l = log.New(os.Stderr, "", log.Lshortfile)
