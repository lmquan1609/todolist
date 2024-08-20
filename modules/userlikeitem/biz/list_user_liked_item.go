package userlikeitembiz

import (
	"context"
	"todolist/common"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
)

type ListUserLikeItemStorage interface {
	ListUsers(ctx context.Context, itemId int, paging *common.Paging, moreKeys ...string) ([]common.SimpleUser, error)
}

type listUserLikeItemBiz struct {
	store ListUserLikeItemStorage
}

func NewListUserLikeItemBiz(store ListUserLikeItemStorage) *listUserLikeItemBiz {
	return &listUserLikeItemBiz{store: store}
}

func (biz *listUserLikeItemBiz) ListUserLikedItem(ctx context.Context, itemId int, paging *common.Paging) ([]common.SimpleUser, error) {
	result, err := biz.store.ListUsers(ctx, itemId, paging, "User")
	if err != nil {
		return nil, common.ErrCannotListEntity(userlikeitemmodel.EntityName, err)
	}

	return result, nil
}
