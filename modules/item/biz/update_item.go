package biz

import (
	"context"
	"todolist/common"
	"todolist/modules/item/model"
)

type UpdateItemStore interface {
	GetItemById(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, data *model.TodoItemUpdate) error
}

type updateItemBiz struct {
	store UpdateItemStore
}

func NewUpdateItemBiz(store UpdateItemStore) *updateItemBiz {
	return &updateItemBiz{store: store}
}

func (biz *updateItemBiz) UpdateItem(ctx context.Context, id int, updatedData *model.TodoItemUpdate) error {
	data, err := biz.store.GetItemById(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	if data.Status == "Deleted" {
		return common.ErrEntityDeleted(model.EntityName, err)
	}

	if err := biz.store.UpdateItem(ctx, map[string]interface{}{"id": id}, updatedData); err != nil {
		return common.ErrCannotDeleteEntity(model.EntityName, err)
	}
	return nil
}
