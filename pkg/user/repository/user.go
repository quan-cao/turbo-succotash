package repository

import "doc-translate-go/pkg/user/entity"

type UserRepository interface {
	Create(u *entity.User) (int, error)
	GetByIsid(isid string) (*entity.User, error)
	Update(u *entity.User) error
	DeleteById(id int) error
}
