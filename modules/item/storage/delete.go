package storage

import (
	"context"
	"todolist/modules/item/model"
)

func (s *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {
	db := s.db

	deletedStatus := "Deleted"
	if err := db.Where(cond).Updates(&model.TodoItemUpdate{Status: &deletedStatus}).Error; err != nil {
		return err
	}
	return nil
}
