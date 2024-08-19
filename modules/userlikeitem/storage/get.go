package userlikeitemstorage

import (
	"context"
	"gorm.io/gorm"
	"todolist/common"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
)

func (s *sqlStore) Find(ctx context.Context, userId, itemId int) (*userlikeitemmodel.Like, error) {
	db := s.db
	var data userlikeitemmodel.Like

	if err := db.Where("user_id = ? AND item_id = ?", userId, itemId).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &data, nil
}
