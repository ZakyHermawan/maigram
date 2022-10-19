package service

import (
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/ZakyHermawan/maigram/internal/repository"
)

type SocialMediaService interface {
	CreateSocialMedia(socialMedia *entity.SocialMedia) error
	GetAllSocialMedia() ([]entity.SocialMedia, error)
	UpdateSocialMedia(socialMedia *entity.SocialMedia) error
	DeleteSocialMedia(socialMedia *entity.SocialMedia) error
}

type socialMedia struct {
	repository repository.SocialMediaRepository
}

func NewSocialMediaService(repository repository.SocialMediaRepository) SocialMediaService {
	return &socialMedia{repository}
}

func (socialSrv *socialMedia) CreateSocialMedia(socialMedia *entity.SocialMedia) error {
	return socialSrv.repository.CreateSocialMedia(socialMedia)
}

func (socialSrv *socialMedia) GetAllSocialMedia() ([]entity.SocialMedia, error) {
	return socialSrv.repository.GetAllSocialMedia()
}

func (socialSrv *socialMedia) UpdateSocialMedia(socialMedia *entity.SocialMedia) error {
	err := socialSrv.repository.UpdateSocialMedia(socialMedia)
	if err != nil {
		return err
	}
	return socialSrv.repository.GetSocialMediaById(socialMedia)
}

func (socialSrv *socialMedia) DeleteSocialMedia(socialMedia *entity.SocialMedia) error {
	return socialSrv.repository.DeleteSocialMedia(socialMedia)
}
