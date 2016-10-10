package dao

import (
	"github.com/jeremyletang/babakoto_api/domain"
	"github.com/jinzhu/gorm"
)

type UserSignupVerification struct {
	db *gorm.DB
}

func NewUserSignupVerificationDao(db *gorm.DB) *UserSignupVerification {
	return &UserSignupVerification{db: db}
}

func (usvd *UserSignupVerification) GetById(id string) (domain.UserSignupVerification, error) {
	usv := domain.UserSignupVerification{Id: id}
	err := usvd.db.First(&usv).Error
	return usv, err
}

func (usvd *UserSignupVerification) GetByUserId(userId string) (domain.UserSignupVerification, error) {
	usv := domain.UserSignupVerification{UserId: userId}
	err := usvd.db.First(&usv).Error
	return usv, err
}

func (usvd *UserSignupVerification) Create(at domain.UserSignupVerification) error {
	return usvd.db.Create(&at).Error
}

func (usvd *UserSignupVerification) Delete(id string) error {
	return usvd.db.Delete(&domain.UserSignupVerification{Id: id}).Error
}
