package dao

import (
	"github.com/jeremyletang/babakoto_api/domain"
	"github.com/jinzhu/gorm"
)

type AccessToken struct {
	db *gorm.DB
}

func NewAccessTokenDao(db *gorm.DB) *AccessToken {
	return &AccessToken{db: db}
}

func (atd *AccessToken) GetById(id string) (domain.AccessToken, error) {
	at := domain.AccessToken{Id: id}
	err := atd.db.First(&at).Error
	return at, err
}

func (atd *AccessToken) GetByUserId(userId string) (domain.AccessToken, error) {
	u := domain.AccessToken{UserId: userId}
	err := atd.db.First(&u).Error
	return u, err
}

func (atd *AccessToken) Create(at domain.AccessToken) error {
	return atd.db.Create(&at).Error
}

func (atd *AccessToken) Delete(id string) error {
	return atd.db.Delete(&domain.AccessToken{Id: id}).Error
}
