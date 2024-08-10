package storage

import (
	"context"
	"todolist/common"
	"todolist/modules/item/model"

	"gorm.io/gorm"
)

func (s *sqlStore) GetItemById(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error) {
	db := s.db
	var data model.TodoItem

	if err := db.Where(cond).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &data, nil
}
