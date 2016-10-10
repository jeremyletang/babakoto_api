package dao

import (
	"github.com/jeremyletang/babakoto_api/domain"
	"github.com/jinzhu/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *User {
	return &User{db: db}
}

func (ud *User) GetById(id string) (domain.User, error) {
	u := domain.User{Id: id}
	err := ud.db.First(&u).Error
	return u, err
}

func (ud *User) GetByEmailOrUsername(str string) (domain.User, error) {
	u := domain.User{}
	err := ud.db.Where("users.email = ? OR users.username = ?", str, str).
		First(&u).Error
	return u, err
}

func (ud *User) GetByMail(email string) (domain.User, error) {
	u := domain.User{Email: email}
	err := ud.db.First(&u).Error
	return u, err
}

func (ud *User) GetByUsername(username string) (domain.User, error) {
	u := domain.User{Username: username}
	err := ud.db.First(&u).Error
	return u, err
}

func (ud *User) Create(u domain.User) error {
	return ud.db.Create(&u).Error
}

func (ud *User) Update(u domain.User) error {
	return ud.db.Save(&u).Error
}
