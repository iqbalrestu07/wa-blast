package services

import (
	"wa-blast/models"
	"wa-blast/repositories"
	"wa-blast/util"
)

// UserService ...
type UserService interface {
	ValidateAuth(key string) (models.CompaniesAuth, error)
}

type userService struct {
	user repositories.UserRepository
}

func (s *userService) ValidateAuth(key string) (models.CompaniesAuth, error) {

	user, err := s.user.GetAuth("secret_key = ? AND active = ?", key, true)

	if user.ID == "" {
		return user, util.NewError("-1002")
	}

	return user, err
}
