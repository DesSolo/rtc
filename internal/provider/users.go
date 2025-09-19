package provider

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/DesSolo/rtc/internal/models"
	"github.com/DesSolo/rtc/internal/storage"
)

// AuthenticateUser ...
func (p *Provider) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	user, err := p.storage.User(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("storage.User: %w", err)
	}

	if !user.IsEnabled {
		slog.DebugContext(ctx, "user is disabled", "username", username)
		return nil, ErrNotFound
	}

	if !isValidPassword(user.PasswordHash, password) {
		slog.DebugContext(ctx, "invalid password", "username", username)
		return nil, ErrNotFound
	}

	return convertUserToModel(user), nil
}

func isValidPassword(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}

//// ChangePassword ...
//func (p *Provider) ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error {
//	// TODO: implement
//	return nil
//}
//
//// ResetPassword ...
//func (p *Provider) ResetPassword(ctx context.Context, username string) error {
//	// TODO: implement
//	return nil
//}

// ListUsers ...
func (p *Provider) ListUsers(ctx context.Context, q string, limit, offset uint64) ([]*models.User, uint64, error) {
	users, total, err := p.storage.Users(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("storage.Users: %w", err)
	}

	return convertUsersToModels(users), total, nil
}

// CreateUser ...
func (p *Provider) CreateUser(ctx context.Context, user *models.User, password string) error {
	passwordHash, err := p.passwordHash(password)
	if err != nil {
		return fmt.Errorf("p.passwordHash: %w", err)
	}

	if err := p.storage.CreateUser(ctx, convertModelToUser(user, passwordHash)); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return ErrAlreadyExists
		}

		return fmt.Errorf("storage.CreateUser: %w", err)
	}

	return nil
}

func (p *Provider) passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	return string(hash), nil
}

// UpdateUserFields ...
type UpdateUserFields struct {
	IsEnabled *bool
	Roles     []string
}

// UpdateUser ...
func (p *Provider) UpdateUser(ctx context.Context, username string, fields *UpdateUserFields) error {
	user, err := p.storage.User(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
	}

	if fields.IsEnabled != nil {
		user.IsEnabled = *fields.IsEnabled
	}

	if fields.Roles != nil {
		user.Roles = fields.Roles
	}

	if err := p.storage.UpdateUser(ctx, user.ID, user); err != nil {
		return fmt.Errorf("p.UpdateUser: %w", err)
	}

	return nil
}
