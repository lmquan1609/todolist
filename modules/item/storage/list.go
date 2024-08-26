package storage

import (
	"context"
	"todolist/common"
	"todolist/modules/item/model"
)

func (s *sqlStore) ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging, moreKeys ...string) ([]model.TodoItem, error) {
	db := s.db
	var data []model.TodoItem

	db = db.Where("status <> ?", "Deleted")

	requester := ctx.Value(common.CurrentUser).(common.Requester)
	db = db.Where("user_id = ?", requester.GetUserId())

	if filter != nil {
		if val := filter.Status; val != "" {
			db = db.Where("status = ?", val)
		}
	}

	if err := db.Table(model.TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if v := paging.FakeCursor; v != "" {
		uid, err := common.FromBase58(v)

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("id < ?", uid.GetLocalID())
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.Limit(paging.Limit).
		Order("id desc").
		Find(&data).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	return data, nil
}
