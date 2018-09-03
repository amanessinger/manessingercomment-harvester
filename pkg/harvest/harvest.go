package harvest

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"strconv"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	Config: aws.Config{
		Region: aws.String(endpoints.UsEast2RegionID),
	},
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = sqs.New(sess)

func FetchComments() (error, []*Comment, []string) {
	out, err := svc.ReceiveMessage(&ReceiveMessageInput)
	if err != nil {
		return err, nil, nil
	}
	var receiptHandles []string
	var comments []*Comment
	for _, m := range (*out).Messages {
		receiptHandles = append(receiptHandles, *(m.ReceiptHandle))
		var c Comment
		json.Unmarshal([]byte(*(m.Body)), &c)
		comments = append(comments, &c)
	}

	return nil, comments, receiptHandles
}

func CleanupQueue(receiptHandles []string) {
	var entries []*sqs.DeleteMessageBatchRequestEntry

	for i, rh := range receiptHandles {
		id := strconv.Itoa(i)
		handle := rh // CAUTION: we need a separate variable, because the address of rh is always the same!
		entries = append(entries, &sqs.DeleteMessageBatchRequestEntry{
			Id:            &id,
			ReceiptHandle: &handle,
		})
	}

	var input = sqs.DeleteMessageBatchInput{
		QueueUrl: &CommentQueueUrl,
		Entries:  entries,
	}
	if out, err := svc.DeleteMessageBatch(&input); err != nil {
		l.Printf("DeleteMessageBatch failed: %v", err)
	} else {
		if len(out.Failed) != 0 {
			for _, failed := range out.Failed {
				id, _ := strconv.Atoi(*((*failed).Id))
				l.Printf("Failed to cleanup %s:\n%s\n%s\n\n",
					*((*failed).Id),
					*((*failed).Message),
					receiptHandles[id])
			}
		}
	}
}
