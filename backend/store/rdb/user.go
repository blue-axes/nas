package rdb

import (
	stdErr "errors"
	"time"

	"github.com/blue-axes/tmpl/pkg/context"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	user struct {
		ID           uint      `gorm:"primarykey"`
		CreatedAt    time.Time `gorm:"autoCreateTime"`
		UpdatedAt    time.Time `gorm:"autoUpdateTime"`
		Username     string    `gorm:"column:username; size:128; uniqueIndex; not null"`
		PasswordHash string    `gorm:"column:password_hash; size:256; not null"`
		CanRead      bool      `gorm:"column:can_read; default: true"`
		CanWrite     bool      `gorm:"column:can_write; default: false"`
		IsAdmin      bool      `gorm:"column:is_admin; default: false"`
	}
)

func (*user) TableName() string {
	return "user"
}

func (m *user) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.PasswordHash = string(hash)
	return nil
}

func (m *user) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(m.PasswordHash), []byte(password)) == nil
}

func (s *txStore) GetUserByUsername(ctx *context.Context, username string) (*user, error) {
	var u user
	err := s.db.Where("username = ?", username).First(&u).Error
	if stdErr.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *txStore) CreateUser(ctx *context.Context, username, password string, canRead, canWrite bool) error {
	u := user{
		Username: username,
		CanRead:  canRead,
		CanWrite: canWrite,
	}
	if err := u.SetPassword(password); err != nil {
		return err
	}
	return s.db.Create(&u).Error
}

func (s *txStore) AdminExists(ctx *context.Context) (bool, error) {
	var count int64
	err := s.db.Model(&user{}).Where("is_admin = ?", true).Count(&count).Error
	return count > 0, err
}

func (s *txStore) CreateAdmin(ctx *context.Context, username, password string) error {
	u := user{
		Username: username,
		CanRead:  true,
		CanWrite: true,
		IsAdmin:  true,
	}
	if err := u.SetPassword(password); err != nil {
		return err
	}
	return s.db.Create(&u).Error
}

func (s *txStore) ListUsers(ctx *context.Context) ([]user, error) {
	var users []user
	err := s.db.Find(&users).Error
	return users, err
}

func (s *txStore) UpdateUserPermissions(ctx *context.Context, username string, canRead, canWrite bool) error {
	return s.db.Model(&user{}).Where("username = ?", username).Updates(map[string]interface{}{
		"can_read":  canRead,
		"can_write": canWrite,
	}).Error
}

func (s *txStore) DeleteUser(ctx *context.Context, username string) error {
	return s.db.Where("username = ?", username).Delete(&user{}).Error
}

func (s *txStore) ChangePassword(ctx *context.Context, username, newPassword string) error {
	var u user
	err := s.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return err
	}
	if err := u.SetPassword(newPassword); err != nil {
		return err
	}
	return s.db.Where("username = ?", username).Updates(map[string]interface{}{
		"password_hash": u.PasswordHash,
	}).Error
}
