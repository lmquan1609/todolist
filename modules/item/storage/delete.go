package storage

import (
	"context"
	"todolist/common"
	"todolist/modules/item/model"

	"gorm.io/gorm"
)

func (s *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {
	db := s.db

	deletedStatus := "Deleted"
	if err := db.Where(cond).Updates(&model.TodoItemUpdate{Status: &deletedStatus}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return common.RecordNotFound
		}
		return common.ErrDB(err)
	}
	return nil
}
