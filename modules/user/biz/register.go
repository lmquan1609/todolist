package userbiz

import (
	"context"
	"todolist/common"
	usermodel "todolist/modules/user/model"
)

type RegisterStorage interface {
	CreateUser(ctx context.Context, data *usermodel.UserCreate) error
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

type Hasher interface {
	Hash(data string) string
}

type registerBiz struct {
	registerStorage RegisterStorage
	hasher          Hasher
}

func NewRegisterBiz(registerStorage RegisterStorage, hasher Hasher) *registerBiz {
	return &registerBiz{
		registerStorage: registerStorage,
		hasher:          hasher,
	}
}

func (biz *registerBiz) Register(ctx context.Context, data *usermodel.UserCreate) error {
	user, err := biz.registerStorage.FindUser(ctx, map[string]interface{}{"email": data.Email})

	if err != nil && err != common.RecordNotFound {
		return common.ErrDB(err)
	}
	if user != nil {
		return usermodel.ErrEmailExisted
	}

	salt := common.GenSalt(50)

	data.Password = biz.hasher.Hash(data.Password + salt)
	data.Salt = salt
	data.Role = usermodel.RoleUser

	if err := biz.registerStorage.CreateUser(ctx, data); err != nil {
		return usermodel.ErrUserNameOrPasswordInvalid
	}
	return nil
}
