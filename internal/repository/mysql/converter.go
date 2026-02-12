package mysql

import (
	"context"

	"github.com/dev-hyunsang/home-library-backend/internal/domain"
	"github.com/dev-hyunsang/home-library-backend/lib/ent"
	"github.com/google/uuid"
)

// UserConverter converts between ent.User and domain.User
type UserConverter struct{}

// ToDomain converts ent.User to domain.User
func (c UserConverter) ToDomain(u *ent.User) *domain.User {
	if u == nil {
		return nil
	}
	return &domain.User{
		ID:            u.ID,
		NickName:      u.NickName,
		Email:         u.Email,
		Password:      u.Password,
		IsPublished:   u.IsPublished,
		IsTermsAgreed: u.IsTermsAgreed,
		FCMToken:      u.FcmToken,
		Timezone:      u.Timezone,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// ToDomainList converts a slice of ent.User to domain.User
func (c UserConverter) ToDomainList(users []*ent.User) []*domain.User {
	result := make([]*domain.User, len(users))
	for i, u := range users {
		result[i] = c.ToDomain(u)
	}
	return result
}

// BookConverter converts between ent.Book and domain.Book
type BookConverter struct{}

// ToDomain converts ent.Book to domain.Book with explicit ownerID
func (c BookConverter) ToDomain(b *ent.Book, ownerID uuid.UUID) *domain.Book {
	if b == nil {
		return nil
	}
	return &domain.Book{
		ID:           b.ID,
		OwnerID:      ownerID,
		Title:        b.BookTitle,
		Author:       b.Author,
		BookISBN:     b.BookIsbn,
		ThumbnailURL: b.ThumbnailURL,
		Status:       b.Status,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}

// ToDomainWithEdges converts ent.Book to domain.Book using loaded edges
func (c BookConverter) ToDomainWithEdges(b *ent.Book) *domain.Book {
	if b == nil {
		return nil
	}

	ownerID := uuid.Nil
	if b.Edges.Owner != nil {
		ownerID = b.Edges.Owner.ID
	} else {
		// Fallback to query if edges not loaded
		ownerID, _ = b.QueryOwner().OnlyID(context.Background())
	}

	return c.ToDomain(b, ownerID)
}

// ToDomainList converts a slice of ent.Book to domain.Book
func (c BookConverter) ToDomainList(books []*ent.Book, ownerID uuid.UUID) []*domain.Book {
	result := make([]*domain.Book, 0, len(books))
	for _, b := range books {
		result = append(result, c.ToDomain(b, ownerID))
	}
	return result
}

// ReviewConverter converts between ent.Review and domain.Review
type ReviewConverter struct{}

// ToDomain converts ent.Review to domain.Review
func (c ReviewConverter) ToDomain(r *ent.Review, ownerID uuid.UUID) *domain.Review {
	if r == nil {
		return nil
	}
	return &domain.Review{
		ID:        r.ID,
		OwnerID:   ownerID,
		BookISBN:  r.BookIsbn,
		Content:   r.Content,
		Rating:    r.Rating,
		IsPublic:  r.IsPublic,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// ToDomainWithEdges converts ent.Review to domain.Review using loaded edges
func (c ReviewConverter) ToDomainWithEdges(r *ent.Review) *domain.Review {
	if r == nil {
		return nil
	}

	ownerID := uuid.Nil
	if r.Edges.Owner != nil {
		ownerID = r.Edges.Owner.ID
	}

	return c.ToDomain(r, ownerID)
}

// ToResponse converts ent.Review to domain.ReviewResponse with owner nickname
func (c ReviewConverter) ToResponse(r *ent.Review) *domain.ReviewResponse {
	if r == nil {
		return nil
	}

	ownerID := uuid.Nil
	nickname := ""
	if r.Edges.Owner != nil {
		ownerID = r.Edges.Owner.ID
		nickname = r.Edges.Owner.NickName
	}

	return &domain.ReviewResponse{
		ID:            r.ID,
		OwnerID:       ownerID,
		OwnerNickname: nickname,
		BookISBN:      r.BookIsbn,
		Content:       r.Content,
		Rating:        r.Rating,
		IsPublic:      r.IsPublic,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

// ReminderConverter converts between ent.ReadingReminder and domain.ReadingReminder
type ReminderConverter struct{}

// ToDomain converts ent.ReadingReminder to domain.ReadingReminder
func (c ReminderConverter) ToDomain(r *ent.ReadingReminder, userID uuid.UUID) *domain.ReadingReminder {
	if r == nil {
		return nil
	}
	return &domain.ReadingReminder{
		ID:           r.ID,
		UserID:       userID,
		ReminderTime: r.ReminderTime,
		DayOfWeek:    domain.DayOfWeek(r.DayOfWeek),
		IsEnabled:    r.IsEnabled,
		Message:      r.Message,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

// ToDomainWithEdges converts ent.ReadingReminder to domain.ReadingReminder using loaded edges
func (c ReminderConverter) ToDomainWithEdges(r *ent.ReadingReminder) *domain.ReadingReminder {
	if r == nil {
		return nil
	}

	userID := uuid.Nil
	if r.Edges.Owner != nil {
		userID = r.Edges.Owner.ID
	}

	return c.ToDomain(r, userID)
}

// EmailVerificationConverter converts ent.EmailVerification
type EmailVerificationConverter struct{}

// ToDomain converts ent.EmailVerification to domain.EmailVerification
func (c EmailVerificationConverter) ToDomain(v *ent.EmailVerification) *domain.EmailVerification {
	if v == nil {
		return nil
	}
	return &domain.EmailVerification{
		ID:         v.ID,
		Email:      v.Email,
		Code:       v.Code,
		ExpiresAt:  v.ExpiresAt,
		IsVerified: v.IsVerified,
		CreatedAt:  v.CreatedAt,
	}
}
