package userstorage

import (
	"context"
	"gorm.io/gorm"
	"todolist/common"
	usermodel "todolist/modules/user/model"
)

func (s *sqlStore) FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*usermodel.User, error) {
	db := s.db
	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	var user usermodel.User

	if err := db.Where(conds).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &user, nil
}
