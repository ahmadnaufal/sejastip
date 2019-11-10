package infra

import (
	"context"
	"encoding/json"
	"log"

	"sejastip.id/api/entity"

	"cloud.google.com/go/pubsub"
)

type PubsubClient struct {
	client *pubsub.Client
}

func NewPubsubClient(projectID string) (*PubsubClient, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &PubsubClient{client}, nil
}

func (p *PubsubClient) PublishNotification(ctx context.Context, notification *entity.NotificationRequest) {
	if p.client == nil {
		return
	}

	b, err := json.Marshal(notification)
	if err != nil {
		return
	}

	topic := p.client.Topic("send-push-notification")
	_, err = topic.Publish(ctx, &pubsub.Message{Data: b}).Get(ctx)
	log.Printf("Published a notification data to Pub/Sub for User ID %d: %v", notification.UserID, err)
}
