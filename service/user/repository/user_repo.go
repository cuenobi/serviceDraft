package repository

import (
	"github.com/cuenobi/serviceDraft/domain"
	"github.com/cuenobi/serviceDraft/service/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	// Read
	Fetch() ([]*entity.User, error)
	FindByID(id uint) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)

	// Create
	Create(user *entity.User) error

	// Update
	Update(user *entity.User) error

	// Delete
	DeleteByID(id uint) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (ur *userRepo) Fetch() ([]*entity.User, error) {
	users := []*entity.User{}
	if err := ur.db.Find(&users).Error; err != nil {
		return nil, domain.ErrNotFound
	}
	return users, nil
}

func (ur *userRepo) FindByID(id uint) (*entity.User, error) {
	user := &entity.User{}
	if err := ur.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (ur *userRepo) FindByUsername(username string) (*entity.User, error) {
	user := &entity.User{}
	if err := ur.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (ur *userRepo) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (ur *userRepo) Create(user *entity.User) error {
	tx := ur.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (ur *userRepo) Update(user *entity.User) error {
	tx := ur.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (ur *userRepo) DeleteByID(id uint) error {
	user := &entity.User{}
	if err := ur.db.Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}
	return nil
}
