package userlikeitemstorage

import (
	"context"
	"todolist/common"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
)

func (s *sqlStore) Create(ctx context.Context, data *userlikeitemmodel.Like) error {
	db := s.db
	if err := db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
