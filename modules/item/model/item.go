package model

import (
	"errors"
	"strings"
	"todolist/common"
)

var (
	ErrTitleCannotBeEmpty = errors.New("title cannot be empty")
	ErrItemIsDeleted      = errors.New("item is deleted")
)

const (
	EntityName = "Item"
)

type TodoItem struct {
	common.SQLModel
	UserId      int                `json:"-" gorm:"column:user_id;"`
	Title       string             `json:"title" gorm:"column:title;"`
	Description string             `json:"description" gorm:"column:description;"`
	Status      string             `json:"status" gorm:"column:status;"`
	Image       *common.Image      `json:"image" gorm:"column:image;"`
	LikedCount  int                `json:"likedCount" gorm:"-"`
	Owner       *common.SimpleUser `json:"owner" gorm:"foreignKey:UserId;"`
}

func (TodoItem) TableName() string {
	return "todo_items"
}

func (i *TodoItem) Mask() {
	i.SQLModel.Mask(common.DBTypeItem)

	if v := i.Owner; i != nil {
		v.Mask()
	}
}

type TodoItemCreation struct {
	Id          int           `json:"id" gorm:"column:id;"`
	UserId      int           `json:"-" gorm:"column:user_id;"`
	Title       string        `json:"title" gorm:"column:title;"`
	Description string        `json:"description" gorm:"column:description;"`
	Image       *common.Image `json:"image" gorm:"column:image;"`
}

func (TodoItemCreation) TableName() string {
	return TodoItem{}.TableName()
}

func (item *TodoItemCreation) Validate() error {
	item.Title = strings.TrimSpace(item.Title)
	if item.Title == "" {
		return ErrTitleCannotBeEmpty
	}
	return nil
}

type TodoItemUpdate struct {
	Id          int     `json:"id" gorm:"column:id;"`
	Description *string `json:"description" gorm:"column:description;"`
	Status      *string `json:"status" gorm:"column:status;"`
}

func (TodoItemUpdate) TableName() string {
	return TodoItem{}.TableName()
}
