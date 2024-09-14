package memcache

import (
	"context"
	"fmt"
	"sync"
	"time"
	usermodel "todolist/modules/user/model"
)

type RealStore interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

type userCaching struct {
	store     Cache
	realStore RealStore
	once      *sync.Once
}

func NewUserCaching(store Cache, realStore RealStore) *userCaching {
	return &userCaching{
		store:     store,
		realStore: realStore,
		once:      new(sync.Once),
	}
}

func (uc *userCaching) FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*usermodel.User, error) {
	var user usermodel.User
	userId := conds["id"].(int)
	key := fmt.Sprintf("user-%d", userId)

	err := uc.store.Get(ctx, key, &user)

	if err == nil && user.Id > 0 {
		return &user, nil
	}

	var userErr error

	uc.once.Do(func() {
		realUser, userErr := uc.realStore.FindUser(ctx, conds, moreInfo...)
		if userErr != nil {
			panic(err)
		}
		user = *realUser
		_ = uc.store.Set(ctx, key, realUser, time.Hour*2)
	})

	if userErr != nil {
		return nil, userErr
	}

	err = uc.store.Get(ctx, key, &user)
	if err == nil && user.Id > 0 {
		return &user, nil
	}
	return nil, err
}
