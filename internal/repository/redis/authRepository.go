package memory

import (
	"fmt"
	"log"

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

	sess.Set("user_id", userID)

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

	raw := sess.Get("user_id")
	if raw == nil {
		return "", fmt.Errorf("user not logged in: %s", userID)
	}

	userID, ok := raw.(string)
	log.Println(userID)

	if !ok {
		return "", fmt.Errorf("invalid user ID type: %T", raw)
	}

	return userID, nil
}

func (repo *AuthRepository) GetAllSession(ctx *fiber.Ctx) ([]string, error) {
	sess, err := repo.store.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return sess.Keys(), nil
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
