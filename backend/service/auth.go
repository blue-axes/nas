package service

import (
	"fmt"

	"github.com/blue-axes/tmpl/pkg/context"
)

type UserInfo struct {
	Username string
	CanRead  bool
	CanWrite bool
	IsAdmin  bool
}

func (svc *Service) ValidateUser(username, password string) (*UserInfo, bool) {
	ctx := context.New()
	u, err := svc.store.RDB().GetUserByUsername(ctx, username)
	if err != nil || u == nil {
		return nil, false
	}
	if u.CheckPassword(password) {
		return &UserInfo{
			Username: u.Username,
			CanRead:  u.CanRead,
			CanWrite: u.CanWrite,
			IsAdmin:  u.IsAdmin,
		}, true
	}
	return nil, false
}

func (svc *Service) InitAdminUser() error {
	ctx := context.New()
	exists, err := svc.store.RDB().AdminExists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return svc.store.RDB().CreateAdmin(ctx, "admin", "admin")
	}
	return nil
}

type ListUserItem struct {
	Username string `json:"Username"`
	CanRead  bool   `json:"CanRead"`
	CanWrite bool   `json:"CanWrite"`
	IsAdmin  bool   `json:"IsAdmin"`
}

func (svc *Service) ListUsers() ([]ListUserItem, error) {
	ctx := context.New()
	users, err := svc.store.RDB().ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]ListUserItem, 0, len(users))
	for _, u := range users {
		result = append(result, ListUserItem{
			Username: u.Username,
			CanRead:  u.CanRead,
			CanWrite: u.CanWrite,
			IsAdmin:  u.IsAdmin,
		})
	}
	return result, nil
}

func (svc *Service) CreateUser(username, password string, canRead, canWrite bool) error {
	ctx := context.New()
	return svc.store.RDB().CreateUser(ctx, username, password, canRead, canWrite)
}

func (svc *Service) UpdateUserPermissions(username string, canRead, canWrite bool) error {
	ctx := context.New()
	return svc.store.RDB().UpdateUserPermissions(ctx, username, canRead, canWrite)
}

func (svc *Service) DeleteUser(username string) error {
	ctx := context.New()
	return svc.store.RDB().DeleteUser(ctx, username)
}

func (svc *Service) ChangePassword(username, currentPassword, newPassword string) error {
	info, ok := svc.ValidateUser(username, currentPassword)
	if !ok {
		return fmt.Errorf("invalid current password")
	}
	_ = info
	ctx := context.New()
	return svc.store.RDB().ChangePassword(ctx, username, newPassword)
}
