package storage

import (
	"context"
	"todolist/modules/item/model"
)

func (s *sqlStore) UpdateItem(ctx context.Context, cond map[string]interface{}, data *model.TodoItemUpdate) error {
	db := s.db
	if err := db.Where(cond).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
