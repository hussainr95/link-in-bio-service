package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hussainr95/link-in-bio-service/internal/entity"
	"github.com/hussainr95/link-in-bio-service/internal/usecase"
	"github.com/stretchr/testify/assert"
)

type mockLinkRepository struct {
	links  map[string]*entity.Link
	nextID int
}

func newMockLinkRepository() *mockLinkRepository {
	return &mockLinkRepository{
		links:  make(map[string]*entity.Link),
		nextID: 1,
	}
}

func (r *mockLinkRepository) Create(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	id := fmt.Sprintf("%d", r.nextID)
	r.nextID++
	link.ID = id
	r.links[id] = link
	return link, nil
}

func (r *mockLinkRepository) GetByID(ctx context.Context, id string) (*entity.Link, error) {
	link, exists := r.links[id]
	if !exists {
		return nil, errors.New("link not found")
	}
	return link, nil
}

func (r *mockLinkRepository) Update(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	if _, exists := r.links[link.ID]; !exists {
		return nil, errors.New("link not found")
	}
	r.links[link.ID] = link
	return link, nil
}

func (r *mockLinkRepository) Delete(ctx context.Context, id string) error {
	if _, exists := r.links[id]; !exists {
		return errors.New("link not found")
	}
	delete(r.links, id)
	return nil
}

func (r *mockLinkRepository) IncrementClicks(ctx context.Context, id string) error {
	link, exists := r.links[id]
	if !exists {
		return errors.New("link not found")
	}
	link.Clicks++
	return nil
}

func (r *mockLinkRepository) DeleteExpired(ctx context.Context) error {
	now := time.Now()
	for id, link := range r.links {
		if link.ExpiresAt.Before(now) {
			delete(r.links, id)
		}
	}
	return nil
}

type mockVisitRepository struct {
	visits []*entity.Visit
}

func newMockVisitRepository() *mockVisitRepository {
	return &mockVisitRepository{
		visits: make([]*entity.Visit, 0),
	}
}

func (r *mockVisitRepository) Create(ctx context.Context, visit *entity.Visit) (*entity.Visit, error) {
	r.visits = append(r.visits, visit)
	return visit, nil
}

// --- Usecase Tests ---

func TestCreateLink(t *testing.T) {
	ctx := context.Background()
	linkRepo := newMockLinkRepository()
	visitRepo := newMockVisitRepository()
	uc := usecase.NewLinkUsecase(linkRepo, visitRepo)

	link := &entity.Link{
		Title:     "Test Link",
		URL:       "http://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	createdLink, err := uc.CreateLink(ctx, link)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdLink.ID)
	assert.Equal(t, 0, createdLink.Clicks)
	assert.False(t, createdLink.CreatedAt.IsZero())
	expectedExpiry := createdLink.CreatedAt.Add(2 * time.Minute)
	diff := createdLink.ExpiresAt.Sub(expectedExpiry)
	assert.LessOrEqual(t, diff.Abs().Seconds(), 1.0, "ExpiresAt should be 2 minutes after CreatedAt")
}

func TestGetLink(t *testing.T) {
	ctx := context.Background()
	linkRepo := newMockLinkRepository()
	visitRepo := newMockVisitRepository()
	uc := usecase.NewLinkUsecase(linkRepo, visitRepo)

	link := &entity.Link{
		Title:     "Test Link",
		URL:       "http://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	createdLink, _ := uc.CreateLink(ctx, link)

	fetchedLink, err := uc.GetLink(ctx, createdLink.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdLink.Title, fetchedLink.Title)
}

func TestUpdateLink(t *testing.T) {
	ctx := context.Background()
	linkRepo := newMockLinkRepository()
	visitRepo := newMockVisitRepository()
	uc := usecase.NewLinkUsecase(linkRepo, visitRepo)

	link := &entity.Link{
		Title:     "Test Link",
		URL:       "http://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	createdLink, _ := uc.CreateLink(ctx, link)

	// Update title and URL
	createdLink.Title = "Updated Title"
	createdLink.URL = "http://updated.com"

	updatedLink, err := uc.UpdateLink(ctx, createdLink)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedLink.Title)
	assert.Equal(t, "http://updated.com", updatedLink.URL)
}

func TestDeleteLink(t *testing.T) {
	ctx := context.Background()
	linkRepo := newMockLinkRepository()
	visitRepo := newMockVisitRepository()
	uc := usecase.NewLinkUsecase(linkRepo, visitRepo)

	link := &entity.Link{
		Title:     "Test Link",
		URL:       "http://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	createdLink, _ := uc.CreateLink(ctx, link)

	err := uc.DeleteLink(ctx, createdLink.ID)
	assert.NoError(t, err)

	_, err = uc.GetLink(ctx, createdLink.ID)
	assert.Error(t, err)
}

func TestVisitLink(t *testing.T) {
	ctx := context.Background()
	linkRepo := newMockLinkRepository()
	visitRepo := newMockVisitRepository()
	uc := usecase.NewLinkUsecase(linkRepo, visitRepo)

	link := &entity.Link{
		Title:     "Test Link",
		URL:       "http://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	createdLink, _ := uc.CreateLink(ctx, link)

	visitedLink, err := uc.VisitLink(ctx, createdLink.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, visitedLink.Clicks)
	assert.Equal(t, 1, len(visitRepo.visits))
}

func TestCleanupExpiredLinks(t *testing.T) {
	ctx := context.Background()
	linkRepo := newMockLinkRepository()
	visitRepo := newMockVisitRepository()
	uc := usecase.NewLinkUsecase(linkRepo, visitRepo)

	// Create an "expired" link using CreateLink, then override ExpiresAt.
	expiredLink := &entity.Link{
		Title: "Expired Link",
		URL:   "http://expired.com",
	}
	createdExpired, _ := uc.CreateLink(ctx, expiredLink)
	// Overwrite ExpiresAt to simulate expiration (1 minute in the past)
	createdExpired.ExpiresAt = time.Now().Add(-1 * time.Minute)

	// Create a valid link; CreateLink sets ExpiresAt to 2 minutes from creation.
	validLink := &entity.Link{
		Title: "Valid Link",
		URL:   "http://valid.com",
	}
	validCreated, _ := uc.CreateLink(ctx, validLink)

	// Run cleanup: this should remove links where ExpiresAt is less than now.
	err := uc.CleanupExpiredLinks(ctx)
	assert.NoError(t, err)

	// The expired link should be removed.
	_, err = uc.GetLink(ctx, createdExpired.ID)
	assert.Error(t, err)

	// The valid link should still exist.
	fetchedValid, err := uc.GetLink(ctx, validCreated.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Valid Link", fetchedValid.Title)
}
