package userlikeitemmodel

import (
	"fmt"
	"time"
	"todolist/common"
)

const EntityName = "UserLikeItem"

type Like struct {
	UserId    int                `json:"user_id" gorm:"column:user_id;"`
	ItemId    int                `json:"item_id" gorm:"column:item_id;"`
	CreatedAt *time.Time         `json:"created_at" gorm:"column:created_at;"`
	User      *common.SimpleUser `json:"-" gorm:"foreignKey:UserId;"`
}

func (Like) TableName() string {
	return "user_like_items"
}

func ErrCannnotLikeItem(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("Cannot like this item"),
		fmt.Sprintf("ErrCannotLikeItem"),
	)
}

func ErrCannnotUnlikeItem(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("Cannot unlike this item"),
		fmt.Sprintf("ErrCannotUnlikeItem"),
	)
}

func ErrDidNotLikeItem(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("You have not liked this item"),
		fmt.Sprintf("ErrDidnotUnlikeItem"),
	)
}
