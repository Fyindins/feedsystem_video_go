package repository

import (
	"feedsystem_video_go/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) RenameByID(id uint, newUsername string) error {
	if err := ur.db.Model(&model.User{}).Where("id = ?", id).Update("username", newUsername).Error; err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) ChangePassword(id uint, newPassword string) error {
	if err := ur.db.Model(&model.User{}).Where("id = ?", id).Update("password", newPassword).Error; err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := ur.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := ur.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
