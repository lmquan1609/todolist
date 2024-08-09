package storage

import (
	"context"
	"todolist/modules/item/model"
)

func (s *sqlStore) GetItemById(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error) {
	db := s.db
	var data model.TodoItem

	if err := db.Where(cond).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
