package userlikeitemstorage

import (
	"context"
	"todolist/common"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
)

func (s *sqlStore) Delete(ctx context.Context, userId, itemId int) error {
	db := s.db

	if err := db.Table(userlikeitemmodel.Like{}.TableName()).
		Where("user_id = ? AND item_id = ?", userId, itemId).
		Delete(nil).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
