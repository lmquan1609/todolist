package subscriber

import (
	"context"
	goservice "github.com/200Lab-Education/go-sdk"
	"log"
	"todolist/pubsub"
)

type HasUserId interface {
	GetUserId() int
}

func PushNotificationAfterUserLikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Push notification after user likes item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			data := message.Data().(HasUserId)
			log.Println("Push notification to user id:", data.GetUserId())
			return nil
		},
	}
}
