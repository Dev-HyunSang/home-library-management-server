package memory

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AuthRepository struct {
	store *session.Store
}

func NewAuthRepository(store *session.Store) *AuthRepository {
	return &AuthRepository{
		store: store,
	}
}

func (repo *AuthRepository) SetSession(userID string, ctx *fiber.Ctx) error {
	sess, err := repo.store.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	sess.Set("user", userID)

	err = sess.Save()
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

func (repo *AuthRepository) GetSessionByID(userID string, ctx *fiber.Ctx) (string, error) {
	sess, err := repo.store.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	raw := sess.Get("user")
	if raw == nil {
		return "", fmt.Errorf("user not logged in: %s", userID)
	}
	userID, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID type: %T", raw)
	}

	return userID, nil
}

func (repo *AuthRepository) DeleteSession(ctx *fiber.Ctx) error {
	sess, err := repo.store.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if err := sess.Destroy(); err != nil {
		return fmt.Errorf("failed to destroy session: %w", err)
	}

	return nil
}
