package middleware

import (
	"context"
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"strings"
	"todolist/common"
	usermodel "todolist/modules/user/model"
	"todolist/plugin/tokenprovider"
)

type AuthenStore interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("Wrong authen header"),
		fmt.Sprintf("ErrWrongAuthHeader"),
	)
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	// "Authorization" : "Bearer {token}"

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader(nil)
	}
	return parts[1], nil
}

// Required Auth
// 1. Get token from header
// 2. Validate token and parse to payload
// 3. From the token payload, we use user_id to find from DB

func RequiredAuth(authStore AuthenStore, serviceCtx goservice.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))

		if err != nil {
			panic(err)
		}

		tokenProvider := serviceCtx.MustGet(common.PluginJWT).(tokenprovider.Provider)

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			panic(err)
		}

		user, err := authStore.FindUser(c.Request.Context(), map[string]interface{}{"id": payload.UserId})
		if err != nil {
			panic(err)
		}

		c.Set(common.CurrentUser, user)
		c.Next()

	}
}
