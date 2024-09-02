package userlikeitembiz

import (
	"context"
	"log"
	"todolist/common/asyncjob"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
)

type UserLikeItemStorage interface {
	Create(ctx context.Context, data *userlikeitemmodel.Like) error
}

type IncreaseItemStorage interface {
	IncreaseLikeCount(ctx context.Context, id int) error
}

type userLikeItemBiz struct {
	store     UserLikeItemStorage
	itemStore IncreaseItemStorage
}

func NewUserLikeItemBiz(store UserLikeItemStorage, itemStore IncreaseItemStorage) *userLikeItemBiz {
	return &userLikeItemBiz{store: store, itemStore: itemStore}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *userlikeitemmodel.Like) error {
	if err := biz.store.Create(ctx, data); err != nil {
		return userlikeitemmodel.ErrCannnotLikeItem(err)
	}

	// go func() {
	// 	if err := biz.itemStore.IncreaseLikeCount(ctx, data.ItemId); err != nil {
	// 		log.Println(err)
	// 	}
	// }()

	job := asyncjob.NewJob(func(ctx context.Context) error {
		if err := biz.itemStore.IncreaseLikeCount(ctx, data.ItemId); err != nil {
			return err
		}
		return nil
	})

	if err := asyncjob.NewGroup(true, job).Run(ctx); err != nil {
		log.Println(err)
	}

	return nil
}
