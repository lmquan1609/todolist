package biz

import (
	"context"
	"todolist/modules/item/model"
)

type GetItemStorage interface {
	GetItemById(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
}

type getItemByIdBiz struct {
	store GetItemStorage
}

func NewGetItemByIdBiz(store GetItemStorage) *getItemByIdBiz {
	return &getItemByIdBiz{store: store}
}

func (biz *getItemByIdBiz) GetItemById(ctx context.Context, id int) (*model.TodoItem, error) {
	data, err := biz.store.GetItemById(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	return data, nil
}
