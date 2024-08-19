package userlikeitembiz

import (
	"context"
	"todolist/common"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
)

type UserUnlikeItemStorage interface {
	Find(ctx context.Context, userId, itemId int) (*userlikeitemmodel.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

type userUnlikeItemBiz struct {
	store UserUnlikeItemStorage
}

func NewUserUnlikeItemBiz(store UserUnlikeItemStorage) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store}
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
	return nil
}
