package userlikeitembiz

import (
	"context"
	"log"
	"todolist/common"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
	"todolist/pubsub"
)

type UserUnlikeItemStorage interface {
	Find(ctx context.Context, userId, itemId int) (*userlikeitemmodel.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

// type DecreaseItemStorage interface {
// 	DecreaseLikeCount(ctx context.Context, id int) error
// }

type userUnlikeItemBiz struct {
	store UserUnlikeItemStorage
	// itemStore DecreaseItemStorage
	ps pubsub.PubSub
}

func NewUserUnlikeItemBiz(store UserUnlikeItemStorage, ps pubsub.PubSub) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store, ps: ps}
}

func (biz *userUnlikeItemBiz) UnlikeItem(ctx context.Context, userId, itemId int) error {
	_, err := biz.store.Find(ctx, userId, itemId)

	if err == common.RecordNotFound {
		return userlikeitemmodel.ErrDidNotLikeItem(err)
	}

	if err != nil {
		return userlikeitemmodel.ErrCannnotUnlikeItem(err)
	}

	if err := biz.store.Delete(ctx, userId, itemId); err != nil {
		return userlikeitemmodel.ErrCannnotUnlikeItem(err)
	}

	// go func() {
	// 	if err := biz.itemStore.DecreaseLikeCount(ctx, itemId); err != nil {
	// 		log.Println(err)
	// 	}
	// }()

	// job := asyncjob.NewJob(func(ctx context.Context) error {
	// 	if err := biz.itemStore.DecreaseLikeCount(ctx, itemId); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })

	// if err := asyncjob.NewGroup(true, job).Run(ctx); err != nil {
	// 	log.Println(err)
	// }

	if err := biz.ps.Publish(ctx, common.TopicUserUnlikedItem, pubsub.NewMessage(&userlikeitemmodel.Like{UserId: userId, ItemId: itemId})); err != nil {
		log.Println(err)
	}

	return nil
}
