package userlikeitembiz

import (
	"context"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
)

type UserLikeItemStorage interface {
	Create(ctx context.Context, data *userlikeitemmodel.Like) error
}

type userLikeItemBiz struct {
	store UserLikeItemStorage
}

func NewUserLikeItemBiz(store UserLikeItemStorage) *userLikeItemBiz {
	return &userLikeItemBiz{store: store}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *userlikeitemmodel.Like) error {
	if err := biz.store.Create(ctx, data); err != nil {
		return userlikeitemmodel.ErrCannnotLikeItem(err)
	}
	return nil
}
