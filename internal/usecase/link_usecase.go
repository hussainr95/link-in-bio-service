package usecase

import (
	"context"
	"time"

	"github.com/hussainr95/link-in-bio-service/internal/entity"
	"github.com/hussainr95/link-in-bio-service/internal/repository"
)

type LinkUsecase interface {
	CreateLink(ctx context.Context, link *entity.Link) (*entity.Link, error)
	GetLink(ctx context.Context, id string) (*entity.Link, error)
	UpdateLink(ctx context.Context, link *entity.Link) (*entity.Link, error)
	DeleteLink(ctx context.Context, id string) error
	VisitLink(ctx context.Context, id string) (*entity.Link, error)
	CleanupExpiredLinks(ctx context.Context) error
}

type linkUsecase struct {
	repo      repository.LinkRepository
	visitRepo repository.VisitRepository
}

func NewLinkUsecase(repo repository.LinkRepository, visitRepo repository.VisitRepository) LinkUsecase {
	return &linkUsecase{repo: repo, visitRepo: visitRepo}
}

func (u *linkUsecase) CreateLink(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	link.CreatedAt = time.Now()
	link.Clicks = 0 // initialize clicks to zero
	return u.repo.Create(ctx, link)
}

func (u *linkUsecase) GetLink(ctx context.Context, id string) (*entity.Link, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *linkUsecase) UpdateLink(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	return u.repo.Update(ctx, link)
}

func (u *linkUsecase) DeleteLink(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *linkUsecase) VisitLink(ctx context.Context, id string) (*entity.Link, error) {
	// Atomically increment clicks using MongoDB's $inc operator.
	if err := u.repo.IncrementClicks(ctx, id); err != nil {
		return nil, err
	}
	// Record the visit for analytics.
	visit := &entity.Visit{
		LinkID:    id,
		VisitedAt: time.Now(),
	}
	if _, err := u.visitRepo.Create(ctx, visit); err != nil {
		return nil, err
	}
	// Return the updated link.
	return u.repo.GetByID(ctx, id)
}

func (u *linkUsecase) CleanupExpiredLinks(ctx context.Context) error {
	return u.repo.DeleteExpired(ctx)
}
