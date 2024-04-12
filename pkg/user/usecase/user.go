package usecase

import (
	"doc-translate-go/pkg/user/entity"
	"doc-translate-go/pkg/user/repository"
)

type UserUseCase struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepository: userRepository}
}

func (uc *UserUseCase) Persist(u *entity.User) (int, error) {
	return uc.userRepository.Create(u)
}

func (uc *UserUseCase) GetByIsid(isid string) (*entity.User, error) {
	return uc.userRepository.GetByIsid(isid)
}
