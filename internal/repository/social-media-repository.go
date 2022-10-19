package repository

import (
	"github.com/ZakyHermawan/maigram/common"
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/jinzhu/gorm"
)

type SocialMediaRepository interface {
	CreateSocialMedia(socialMedia *entity.SocialMedia) error
	GetAllSocialMedia() ([]entity.SocialMedia, error)
	GetSocialMediaById(socialMedia *entity.SocialMedia) error
	UpdateSocialMedia(socialMedia *entity.SocialMedia) error
	DeleteSocialMedia(socialMedia *entity.SocialMedia) error
}

type socialMediaRepository struct {
	*gorm.DB
}

func NewSocialMediaRepository() SocialMediaRepository {
	return &socialMediaRepository{common.GetDB()}
}

func (repository *socialMediaRepository) CreateSocialMedia(socialMedia *entity.SocialMedia) error {
	return repository.Model(socialMedia).Create(socialMedia).Error
}

func (repository *socialMediaRepository) GetAllSocialMedia() ([]entity.SocialMedia, error) {
	var socialMedias []entity.SocialMedia
	return socialMedias, repository.Model(&entity.SocialMedia{}).Find(&socialMedias).Error
}

func (repository *socialMediaRepository) GetSocialMediaById(socialMedia *entity.SocialMedia) error {
	return repository.Model(socialMedia).Where("id = ?", socialMedia.Model.ID).First(socialMedia).Error
}

func (repository *socialMediaRepository) UpdateSocialMedia(socialMedia *entity.SocialMedia) error {
	return repository.Model(socialMedia).Updates(socialMedia).Error
}

func (repository *socialMediaRepository) DeleteSocialMedia(socialMedia *entity.SocialMedia) error {
	return repository.Model(socialMedia).Delete(socialMedia).Error
}
