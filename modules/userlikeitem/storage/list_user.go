package userlikeitemstorage

import (
	"context"
	"github.com/btcsuite/btcutil/base58"
	"time"
	"todolist/common"
	userlikeitemmodel "todolist/modules/userlikeitem/model"
)

const timeLayout = "2006-01-02T15:04:05.999999"

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

	if v := paging.FakeCursor; v != "" {
		timeCreated, err := time.Parse(timeLayout, string(base58.Decode(v)))

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("created_at < ?", timeCreated.Format(timeLayout))
	} else {
		db = db.Offset(paging.Limit * (paging.Page - 1))
	}

	if err := db.Limit(paging.Limit).
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

func (s *sqlStore) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	result := make(map[int]int)
	db := s.db

	type sqlData struct {
		ItemId int `gorm:"column:item_id"`
		Count  int `gorm:"column:count"`
	}

	var listLike []sqlData

	if err := db.Table(userlikeitemmodel.Like{}.TableName()).
		Select("item_id, COUNT(item_id) as `count`").
		Where("item_id in (?)", ids).
		Group("item_id").Find(&listLike).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for _, item := range listLike {
		result[item.ItemId] = item.Count
	}
	return result, nil
}
