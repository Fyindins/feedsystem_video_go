package account

import (
	"context"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (ar *AccountRepository) CreateAccount(ctx context.Context, account *Account) error {
	if err := ar.db.WithContext(ctx).Create(account).Error; err != nil {
		return err
	}
	return nil
}

func (ar *AccountRepository) Rename(ctx context.Context, id uint, newUsername string) error {
	if err := ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("username", newUsername).Error; err != nil {
		return err
	}
	return nil
}

func (ar *AccountRepository) ChangePassword(ctx context.Context, id uint, newPassword string) error {
	if err := ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("password", newPassword).Error; err != nil {
		return err
	}
	return nil
}

func (ar *AccountRepository) FindByID(ctx context.Context, id uint) (*Account, error) {
	var account Account
	if err := ar.db.WithContext(ctx).First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (ar *AccountRepository) FindByUsername(ctx context.Context, username string) (*Account, error) {
	var account Account
	if err := ar.db.WithContext(ctx).Where("username = ?", username).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (ar *AccountRepository) Login(ctx context.Context, id uint, token string) error {
	if err := ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("token", token).Error; err != nil {
		return err
	}
	return nil
}

func (ar *AccountRepository) Logout(ctx context.Context, id uint) error {
	if err := ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("token", "").Error; err != nil {
		return err
	}
	return nil
}
