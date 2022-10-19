package repository

import (
	"github.com/ZakyHermawan/maigram/common"
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type PhotoRepository interface {
	CreatePhoto(photo *entity.Photo) error
	GetAllPhoto() ([]entity.Photo, error)
	GetPhotoById(photo *entity.Photo) error
	UpdatePhoto(photo *entity.Photo) error
	DeletePhotoById(photo *entity.Photo) error
}

type photoRepository struct {
	*gorm.DB
}

func NewPhotoRepository() PhotoRepository {
	return &photoRepository{common.GetDB()}
}

func (repo *photoRepository) CreatePhoto(photo *entity.Photo) error {
	return repo.Model(photo).Create(photo).Error
}

func (repo *photoRepository) GetAllPhoto() ([]entity.Photo, error) {
	var photos []entity.Photo
	return photos, repo.Model(&entity.Photo{}).Find(&photos).Error
}

func (repo *photoRepository) GetPhotoById(photo *entity.Photo) error {
	return repo.Model(photo).Where("id = ?", photo.Model.ID).First(photo).Error
}

func (repo *photoRepository) UpdatePhoto(photo *entity.Photo) error {
	return repo.Model(photo).Updates(photo).Error
}

func (repo *photoRepository) DeletePhotoById(photo *entity.Photo) error {
	return repo.Model(photo).Delete(photo).Error
}
