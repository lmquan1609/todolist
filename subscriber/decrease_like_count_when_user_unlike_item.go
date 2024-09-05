package subscriber

import (
	"context"
	goservice "github.com/200Lab-Education/go-sdk"
	"gorm.io/gorm"
	"todolist/common"
	"todolist/modules/item/storage"
	"todolist/pubsub"
)

func DecreaseLikeCountWhenUserUnlikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Decrease like count after user unlikes item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
			data := message.Data().(HasItemId)
			return storage.NewSQLStore(db).DecreaseLikeCount(ctx, data.GetItemId())
		},
	}
}
