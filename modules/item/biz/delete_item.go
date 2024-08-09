package biz

import (
	"context"
	"todolist/modules/item/model"
)

type DeleteItemStore interface {
	GetItemById(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	DeleteItem(ctx context.Context, cond map[string]interface{}) error
}

type deleteItemBiz struct {
	store DeleteItemStore
}

func NewDeleteItemBiz(store DeleteItemStore) *deleteItemBiz {
	return &deleteItemBiz{store: store}
}

func (biz *deleteItemBiz) DeleteItemBiz(ctx context.Context, id int) error {
	data, err := biz.store.GetItemById(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	if data.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}

	if err := biz.store.DeleteItem(ctx, map[string]interface{}{"id": id}); err != nil {
		return err
	}
	return nil
}
