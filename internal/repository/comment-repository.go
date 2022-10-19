package repository

import (
	"github.com/ZakyHermawan/maigram/common"
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/jinzhu/gorm"
)

type CommentRepository interface {
	CreateComment(comment *entity.Comment) error
	GetAllComment() ([]entity.Comment, error)
	GetCommentById(comment *entity.Comment) error
	UpdateComment(comment *entity.Comment) error
	DeleteComment(comment *entity.Comment) error
}

type commentRepository struct {
	*gorm.DB
}

func NewCommentRepository() CommentRepository {
	return &commentRepository{common.GetDB()}
}

func (repository *commentRepository) CreateComment(comment *entity.Comment) error {
	return repository.Model(comment).Create(comment).Error
}

func (repository *commentRepository) GetAllComment() ([]entity.Comment, error) {
	var comments []entity.Comment
	return comments, repository.Model(&entity.Comment{}).Find(&comments).Error
}

func (repository *commentRepository) GetCommentById(comment *entity.Comment) error {
	return repository.Model(comment).Where("id = ?", comment.Model.ID).First(comment).Error
}

func (repository *commentRepository) UpdateComment(comment *entity.Comment) error {
	return repository.Model(comment).Updates(comment).Error
}

func (repository *commentRepository) DeleteComment(comment *entity.Comment) error {
	return repository.Model(comment).Delete(comment).Error
}
