package service

import (
	"encoding/hex"
	"errors"
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/ZakyHermawan/maigram/internal/repository"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *entity.User) error
	Login(user *entity.User) (bool, error)
	UpdateUsernameOrEmail(user *entity.User) error
	DeleteUser(user *entity.User) error
}
type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}

func (userSrv *userService) Register(user *entity.User) error {
	err := userSrv.repository.GetUserByEmail(user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			byteArray, encryptErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if encryptErr != nil {
				return err
			}
			user.Password = hex.EncodeToString(byteArray)
			err = userSrv.repository.CreateUser(user)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return errors.New("user already exist")
}

func (userSrv *userService) Login(user *entity.User) (bool, error) {
	passwordFromRequest := user.Password
	err := userSrv.repository.GetUserByEmail(user)
	if err != nil {
		return false, err
	}
	hashedPasswordFromDB, errDecodeFromDB := hex.DecodeString(user.Password)
	if errDecodeFromDB != nil {
		return false, errDecodeFromDB
	}
	err = bcrypt.CompareHashAndPassword(hashedPasswordFromDB, []byte(passwordFromRequest))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (userSrv *userService) UpdateUsernameOrEmail(user *entity.User) error {
	err := userSrv.repository.UpdateUser(user)
	if err != nil {
		return err
	}

	// get updated user
	return userSrv.repository.GetUserByID(user)
}

func (userSrv *userService) DeleteUser(user *entity.User) error {
	err := userSrv.repository.GetUserByEmail(user)
	if err != nil {
		return err
	}
	return userSrv.repository.DeleteUser(user)
}
