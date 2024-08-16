package biz

import (
	"context"
	"todolist/common"
	"todolist/modules/item/model"
)

type ListItemStore interface {
	ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging, moreKeys ...string) ([]model.TodoItem, error)
}

type listItemBiz struct {
	store     ListItemStore
	requester common.Requester
}

func NewListItemBiz(store ListItemStore, requester common.Requester) *listItemBiz {
	return &listItemBiz{store: store, requester: requester}
}

func (biz *listItemBiz) ListItemBiz(ctx context.Context, filter *model.Filter, paging *common.Paging) ([]model.TodoItem, error) {
	paging.Process()

	ctxStore := context.WithValue(ctx, common.CurrentUser, biz.requester)

	data, err := biz.store.ListItem(ctxStore, filter, paging, "Owner")
	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}
	return data, nil
}
