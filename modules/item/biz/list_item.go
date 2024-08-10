package biz

import (
	"context"
	"todolist/common"
	"todolist/modules/item/model"
)

type ListItemStore interface {
	ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging) ([]model.TodoItem, error)
}

type listItemBiz struct {
	store ListItemStore
}

func NewListItemBiz(store ListItemStore) *listItemBiz {
	return &listItemBiz{store: store}
}

func (biz *listItemBiz) ListItemBiz(ctx context.Context, filter *model.Filter, paging *common.Paging) ([]model.TodoItem, error) {
	paging.Process()

	data, err := biz.store.ListItem(ctx, filter, paging)
	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}
	return data, nil
}
