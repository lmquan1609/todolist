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

type TodoItem struct {
	common.SQLModel
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status;"`
}

func (TodoItem) TableName() string {
	return "todo_items"
}

type TodoItemCreation struct {
	Id          int    `json:"id" gorm:"column:id;"`
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
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
