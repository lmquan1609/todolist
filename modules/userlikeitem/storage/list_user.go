package userlikeitemstorage

import (
	"context"
	"todolist/common"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
)

func (s *sqlStore) ListUsers(ctx context.Context, itemId int, paging *common.Paging, moreKeys ...string) ([]common.SimpleUser, error) {
	var data []userlikeitemmodel.Like
	db := s.db

	db = db.Where("item_id = ?", itemId)
	if err := db.Table(userlikeitemmodel.Like{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if err := db.Offset(paging.Limit * (paging.Page - 1)).
		Limit(paging.Limit).
		Order("created_at desc").
		Find(&data).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	result := make([]common.SimpleUser, len(data))

	for i := range data {
		result[i] = *data[i].User
		result[i].UpdatedAt = nil
		result[i].CreatedAt = data[i].CreatedAt
	}
	return result, nil
}
