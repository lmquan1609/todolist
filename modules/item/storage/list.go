package storage

import (
	"context"
	"todolist/common"
	"todolist/modules/item/model"
)

func (s *sqlStore) ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging) ([]model.TodoItem, error) {
	db := s.db
	var data []model.TodoItem

	db = db.Where("status <> ?", "Deleted")
	if filter != nil {
		if val := filter.Status; val != "" {
			db = db.Where("status = ?", val)
		}
	}

	if err := db.Table(model.TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	if err := db.Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Order("id desc").
		Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
