package userlikeitembiz

import (
	"context"
	"log"
	"todolist/common"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
	"todolist/pubsub"
)

type UserLikeItemStorage interface {
	Create(ctx context.Context, data *userlikeitemmodel.Like) error
}

// type IncreaseItemStorage interface {
// 	IncreaseLikeCount(ctx context.Context, id int) error
// }

type userLikeItemBiz struct {
	store UserLikeItemStorage
	// itemStore IncreaseItemStorage
	ps pubsub.PubSub
}

func NewUserLikeItemBiz(store UserLikeItemStorage, ps pubsub.PubSub) *userLikeItemBiz {
	return &userLikeItemBiz{store: store, ps: ps}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *userlikeitemmodel.Like) error {
	if err := biz.store.Create(ctx, data); err != nil {
		return userlikeitemmodel.ErrCannnotLikeItem(err)
	}

	if err := biz.ps.Publish(ctx, common.TopicUserLikedItem, pubsub.NewMessage(data)); err != nil {
		log.Println(err)
	}

	// go func() {
	// 	if err := biz.itemStore.IncreaseLikeCount(ctx, data.ItemId); err != nil {
	// 		log.Println(err)
	// 	}
	// }()

	// job := asyncjob.NewJob(func(ctx context.Context) error {
	// 	if err := biz.itemStore.IncreaseLikeCount(ctx, data.ItemId); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })

	// if err := asyncjob.NewGroup(true, job).Run(ctx); err != nil {
	// 	log.Println(err)
	// }

	return nil
}
