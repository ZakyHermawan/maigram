package service

import (
	"github.com/ZakyHermawan/maigram/entity"
	"github.com/ZakyHermawan/maigram/internal/repository"
)

type CommentService interface {
	CreateComment(comment *entity.Comment) error
	GetAllComment() ([]entity.Comment, error)
	UpdateComment(comment *entity.Comment) error
	DeleteComment(comment *entity.Comment) error
}

type commentService struct {
	repository repository.CommentRepository
}

func NewCommentService(repository repository.CommentRepository) CommentService {
	return &commentService{repository}
}

func (commentSrv *commentService) CreateComment(comment *entity.Comment) error {
	return commentSrv.repository.CreateComment(comment)
}

func (commentSrv *commentService) GetAllComment() ([]entity.Comment, error) {
	return commentSrv.repository.GetAllComment()
}

func (commentSrv *commentService) UpdateComment(comment *entity.Comment) error {
	err := commentSrv.repository.UpdateComment(comment)
	if err != nil {
		return err
	}
	return commentSrv.repository.GetCommentById(comment)
}

func (commentSrv *commentService) DeleteComment(comment *entity.Comment) error {
	return commentSrv.repository.DeleteComment(comment)
}
