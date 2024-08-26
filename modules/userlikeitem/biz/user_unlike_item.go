package userlikeitembiz

import (
	"context"
	"log"
	"todolist/common"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
)

type UserUnlikeItemStorage interface {
	Find(ctx context.Context, userId, itemId int) (*userlikeitemmodel.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

type DecreaseItemStorage interface {
	DecreaseLikeCount(ctx context.Context, id int) error
}

type userUnlikeItemBiz struct {
	store     UserUnlikeItemStorage
	itemStore DecreaseItemStorage
}

func NewUserUnlikeItemBiz(store UserUnlikeItemStorage, itemStore DecreaseItemStorage) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store, itemStore: itemStore}
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

	go func() {
		if err := biz.itemStore.DecreaseLikeCount(ctx, itemId); err != nil {
			log.Println(err)
		}
	}()

	return nil
}
