package biz

import (
	"context"
	"errors"
	"todolist/common"
	"todolist/modules/item/model"
)

type DeleteItemStore interface {
	GetItemById(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	DeleteItem(ctx context.Context, cond map[string]interface{}) error
}

type deleteItemBiz struct {
	store     DeleteItemStore
	requester common.Requester
}

func NewDeleteItemBiz(store DeleteItemStore, requester common.Requester) *deleteItemBiz {
	return &deleteItemBiz{store: store, requester: requester}
}

func (biz *deleteItemBiz) DeleteItemBiz(ctx context.Context, id int) error {
	data, err := biz.store.GetItemById(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	if data.Status == "Deleted" {
		return common.ErrEntityDeleted(model.EntityName, err)
	}

	if biz.requester.GetUserId() != data.UserId {
		return common.ErrNoPermission(errors.New("No permission"))
	}

	if err := biz.store.DeleteItem(ctx, map[string]interface{}{"id": id}); err != nil {
		return common.ErrCannotDeleteEntity(model.EntityName, err)
	}
	return nil
}
