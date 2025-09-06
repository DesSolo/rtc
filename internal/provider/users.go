package provider

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"rtc/internal/models"
	"rtc/internal/storage"
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
