package database

import (
	"testing"

	"github.com/Mazzael/go-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})

	user, _ := entity.NewUser("John Doe", "johndoe@example.com", "123456")
	gormUserRepository := NewUser(db)

	err = gormUserRepository.Create(user)
	assert.Nil(t, err)

	var userFound entity.User
	err = db.First(&userFound, "email = ?", user.Email).Error
	assert.Nil(t, err)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotEmpty(t, userFound.Password)
}

func TestFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})

	user, _ := entity.NewUser("John Doe", "johndoe@example.com", "123456")
	gormUserRepository := NewUser(db)

	err = gormUserRepository.Create(user)
	assert.Nil(t, err)

	userFound, err := gormUserRepository.FindByEmail("johndoe@example.com")
	assert.Nil(t, err)
	assert.NotNil(t, userFound)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotEmpty(t, userFound.Password)
}
