package service

import (
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/ZakyHermawan/maigram/internal/repository"
)

type PhotoService interface {
	CreatePhoto(photo *entity.Photo) error
	GetAllPhoto() ([]entity.Photo, error)
	UpdatePhoto(photo *entity.Photo) error
	DeletePhoto(photo *entity.Photo) error
}

type photoService struct {
	repository repository.PhotoRepository
}

func NewPhotoService(repository repository.PhotoRepository) PhotoService {
	return &photoService{repository}
}

func (photoSrv *photoService) CreatePhoto(photo *entity.Photo) error {
	return photoSrv.repository.CreatePhoto(photo)
}

func (photoSrv *photoService) GetAllPhoto() ([]entity.Photo, error) {
	return photoSrv.repository.GetAllPhoto()
}

func (photoSrv *photoService) UpdatePhoto(photo *entity.Photo) error {
	err := photoSrv.repository.UpdatePhoto(photo)
	if err != nil {
		return err
	}
	return photoSrv.repository.GetPhotoById(photo)
}

func (photoSrv *photoService) DeletePhoto(photo *entity.Photo) error {
	return photoSrv.repository.DeletePhotoById(photo)
}
