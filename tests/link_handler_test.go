package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	httphandler "github.com/hussainr95/link-in-bio-service/internal/delivery/http"
	"github.com/hussainr95/link-in-bio-service/internal/entity"
	"github.com/hussainr95/link-in-bio-service/internal/usecase"
	"github.com/stretchr/testify/assert"
)

// Reuse the same mock repositories from link_usecase_test.go

func setupRouter() (*gin.Engine, usecase.LinkUsecase) {
	linkRepo := newFakeLinkRepository()
	visitRepo := newFakeVisitRepository()
	uc := usecase.NewLinkUsecase(linkRepo, visitRepo)
	handler := httphandler.NewLinkHandler(uc)

	router := gin.Default()
	// Register routes (without auth for simplicity)
	router.POST("/links", handler.CreateLink)
	router.GET("/links/:id", handler.GetLink)
	router.PUT("/links/:id", handler.UpdateLink)
	router.DELETE("/links/:id", handler.DeleteLink)
	router.GET("/visit/:id", handler.VisitLink)
	return router, uc
}

func TestCreateLinkEndpoint(t *testing.T) {
	router, _ := setupRouter()

	newLink := entity.Link{
		Title:     "Test Link",
		URL:       "http://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	body, _ := json.Marshal(newLink)

	req, _ := http.NewRequest("POST", "/links", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created entity.Link
	err := json.Unmarshal(w.Body.Bytes(), &created)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, newLink.Title, created.Title)
}

func TestGetLinkEndpoint(t *testing.T) {
	router, uc := setupRouter()
	ctx := context.Background()

	// Create a link directly via the usecase
	link := &entity.Link{
		Title:     "Test Link",
		URL:       "http://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	createdLink, _ := uc.CreateLink(ctx, link)

	req, _ := http.NewRequest("GET", "/links/"+createdLink.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var fetched entity.Link
	err := json.Unmarshal(w.Body.Bytes(), &fetched)
	assert.NoError(t, err)
	assert.Equal(t, createdLink.Title, fetched.Title)
}

func TestUpdateLinkEndpoint(t *testing.T) {
	router, uc := setupRouter()
	ctx := context.Background()

	// Create a link
	link := &entity.Link{
		Title:     "Original Title",
		URL:       "http://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	createdLink, _ := uc.CreateLink(ctx, link)

	// Update via HTTP PUT
	updatedData := entity.Link{
		Title:     "Updated Title",
		URL:       "http://updated.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	body, _ := json.Marshal(updatedData)
	req, _ := http.NewRequest("PUT", "/links/"+createdLink.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var updated entity.Link
	err := json.Unmarshal(w.Body.Bytes(), &updated)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, "http://updated.com", updated.URL)
}

func TestDeleteLinkEndpoint(t *testing.T) {
	router, uc := setupRouter()
	ctx := context.Background()

	// Create a link
	link := &entity.Link{
		Title:     "Test Link",
		URL:       "http://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	createdLink, _ := uc.CreateLink(ctx, link)

	req, _ := http.NewRequest("DELETE", "/links/"+createdLink.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Try fetching the deleted link
	req, _ = http.NewRequest("GET", "/links/"+createdLink.ID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVisitLinkEndpoint(t *testing.T) {
	router, uc := setupRouter()
	ctx := context.Background()

	// Create a link
	link := &entity.Link{
		Title:     "Test Link",
		URL:       "http://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	createdLink, _ := uc.CreateLink(ctx, link)

	req, _ := http.NewRequest("GET", "/visit/"+createdLink.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var visited entity.Link
	err := json.Unmarshal(w.Body.Bytes(), &visited)
	assert.NoError(t, err)
	assert.Equal(t, 1, visited.Clicks)
}
