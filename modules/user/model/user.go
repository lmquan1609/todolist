package usermodel

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"todolist/common"
)

const EntityName = "User"

type UserRole int

const (
	RoleUser UserRole = 1 << iota
	RoleAdmin
	RoleShipper
	RoleMod
)

func (role UserRole) String() string {
	switch role {
	case RoleAdmin:
		return "admin"
	case RoleShipper:
		return "shipper"
	case RoleMod:
		return "mod"
	default:
		return "user"
	}
}

func (role *UserRole) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("Failed to unmarshal JSONB value:", value))
	}
	var r UserRole

	roleVal := string(bytes)
	if roleVal == "user" {
		r = RoleUser
	} else if roleVal == "admin" {
		r = RoleAdmin
	}
	*role = r
	return nil
}

func (role *UserRole) Value() (driver.Value, error) {
	if role == nil {
		return nil, nil
	}
	return role.String(), nil
}

func (role *UserRole) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", role.String())), nil
}

type User struct {
	common.SQLModel
	Email     string   `json:"email" gorm:"column:email;"`
	Password  string   `json:"-" gorm:"column:password;"`
	Salt      string   `json:"-" gorm:"column:salt;"`
	LastName  string   `json:"last_name" gorm:"column:last_name;"`
	FirstName string   `json:"first_name" gorm:"column:first_name;"`
	Phone     string   `json:"phone" gorm:"column:phone;"`
	Role      UserRole `json:"role" gorm:"column:role;"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) GetUserId() int {
	return u.Id
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	return u.Role.String()
}

type UserCreate struct {
	common.SQLModel
	Email     string   `json:"email" gorm:"column:email;"`
	Password  string   `json:"password" gorm:"column:password;"`
	Salt      string   `json:"-" gorm:"column:salt;"`
	LastName  string   `json:"last_name" gorm:"column:last_name;"`
	FirstName string   `json:"first_name" gorm:"column:first_name;"`
	Role      UserRole `json:"-" gorm:"column:role;"`
}

func (UserCreate) TableName() string {
	return User{}.TableName()
}

type UserLogin struct {
	Email    string `json:"email" form:"email" gorm:"column:email;"`
	Password string `json:"password" form:"password" gorm:"column:password;"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}

var (
	ErrUserNameOrPasswordInvalid = common.NewCustomError(
		errors.New("Username or password invalid"),
		"Username or password invalid",
		"ErrUsernameOrPasswordInvalid",
	)

	ErrEmailExisted = common.NewCustomError(
		errors.New("Email has already existed"),
		"Email has already existed",
		"ErrEmailExisted",
	)
)
