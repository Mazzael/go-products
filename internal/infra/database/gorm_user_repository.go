package database

import (
	"github.com/Mazzael/go-api/internal/entity"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{DB: db}
}

func (u *GormUserRepository) Create(user *entity.User) error {
	return u.DB.Create(user).Error
}
func (u *GormUserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
