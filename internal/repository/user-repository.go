package repository

import (
	"github.com/ZakyHermawan/maigram/common"
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByID(user *entity.User) error
	GetUserByEmail(user *entity.User) error
	UpdateUser(user *entity.User) error
	DeleteUser(user *entity.User) error
}

type userRepository struct {
	*gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{common.GetDB()}
}

func (repo *userRepository) CreateUser(user *entity.User) error {
	return repo.DB.Model(user).Create(user).Error
}

func (repo *userRepository) GetUserByID(user *entity.User) error {
	return repo.DB.Model(user).Where("id = ?", user.Model.ID).First(user).Error
}

func (repo *userRepository) GetUserByEmail(user *entity.User) error {
	return repo.DB.Model(user).Where("email = ?", user.Email).First(user).Error
}

func (repo *userRepository) UpdateUser(user *entity.User) error {
	return repo.DB.Model(user).Updates(user).Error
}

func (repo *userRepository) DeleteUser(user *entity.User) error {
	return repo.DB.Model(user).Delete(user).Error
}
