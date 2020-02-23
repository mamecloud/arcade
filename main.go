package main
import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
)

func pullMsgs(w io.Writer, projectID string, subID string, topicID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	// topic of type https://godoc.org/cloud.google.com/go/pubsub#Topic
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
			return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	//topic := client.Topic(topicID)

	// Publish 10 messages on the topic.
	// var results []*pubsub.PublishResult
	// for i := 0; i < 10; i++ {
	// 		res := topic.Publish(ctx, &pubsub.Message{
	// 				Data: []byte(fmt.Sprintf("hello world #%d", i)),
	// 		})
	// 		results = append(results, res)
	// 		fmt.Printf("Published message %d\n", i)
	// }

	// Check that all messages were published.
	// for _, r := range results {
	// 		_, err := r.Get(ctx)
	// 		if err != nil {
	// 				return fmt.Errorf("Get: %v", err)
	// 		}
	// }

	// Consume 10 messages.
	var mu sync.Mutex
	received := 0
	sub := client.Subscription(subID)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
			fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
			msg.Ack()
			mu.Lock()
			defer mu.Unlock()
			received++
			if received == 10 {
					cancel()
			}
	})
	if err != nil {
			return fmt.Errorf("Receive: %v", err)
	}
	return nil
}

func main() {
	projectID := "mamecloud"
	subID := "arcade"
	topicID := "game-request"
	err := pullMsgs(os.Stderr, projectID, subID, topicID)
	fmt.Println(err)
	fmt.Println("Fin.")
}