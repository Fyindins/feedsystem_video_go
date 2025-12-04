package service

import (
	"feedsystem_video_go/internal/model"
	"feedsystem_video_go/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) CreateUser(user *model.User) error {
	if err := us.userRepository.CreateUser(user); err != nil {
		return err
	}
	return nil
}

func (us *UserService) RenameByID(id uint, newUsername string) error {
	if err := us.userRepository.RenameByID(id, newUsername); err != nil {
		return err
	}
	return nil
}

func (us *UserService) ChangePassword(id uint, newPassword string) error {
	if err := us.userRepository.ChangePassword(id, newPassword); err != nil {
		return err
	}
	return nil
}

func (us *UserService) FindByID(id uint) (*model.User, error) {
	if user, err := us.userRepository.FindByID(id); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (us *UserService) FindByUsername(username string) (*model.User, error) {
	if user, err := us.userRepository.FindByUsername(username); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}
