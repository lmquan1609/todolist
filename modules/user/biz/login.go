package userbiz

import (
	"context"
	"todolist/common"
	usermodel "todolist/modules/user/model"
	"todolist/plugin/tokenprovider"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

type loginBiz struct {
	storeUser     LoginStorage
	tokenProvider tokenprovider.Provider
	hasher        Hasher
	expiry        int
}

func NewLoginBiz(storeUser LoginStorage, tokenProvider tokenprovider.Provider, hasher Hasher, expiry int) *loginBiz {
	return &loginBiz{
		storeUser:     storeUser,
		tokenProvider: tokenProvider,
		hasher:        hasher,
		expiry:        expiry,
	}
}

// 1. Find user, email
// 2. Has pass from input and compare with pass in db
// 3. Issue JWT token for client
// 3.1. Access token and refresh token
// 4. Return token(s)

func (biz *loginBiz) Login(ctx context.Context, data *usermodel.UserLogin) (*tokenprovider.Token, error) {
	user, err := biz.storeUser.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if err != nil {
		return nil, usermodel.ErrUserNameOrPasswordInvalid
	}

	passHashed := biz.hasher.Hash(data.Password + user.Salt)

	if user.Password != passHashed {
		return nil, usermodel.ErrUserNameOrPasswordInvalid
	}

	payload := tokenprovider.TokenPayload{
		UserId: user.Id,
		Role:   user.Role.String(),
	}

	accessToken, err := biz.tokenProvider.Generate(payload, biz.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}
	return accessToken, nil
}
