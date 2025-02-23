package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hussainr95/link-in-bio-service/internal/entity"
	"github.com/hussainr95/link-in-bio-service/internal/usecase"
)

type LinkHandler struct {
	usecase usecase.LinkUsecase
}

func NewLinkHandler(u usecase.LinkUsecase) *LinkHandler {
	return &LinkHandler{usecase: u}
}

// RegisterAPIRoutes sets up the routing for link-related endpoints
func (h *LinkHandler) RegisterAPIRoutes(router *gin.Engine) {
	router.POST("/links", h.CreateLink)
	router.GET("/links/:id", h.GetLink)
	router.PUT("/links/:id", h.UpdateLink)
	router.DELETE("/links/:id", h.DeleteLink)
	router.GET("/visit/:id", h.VisitLink)
}

// CreateLink handles POST /links
// CreateLink godoc
// @Summary Create a new link
// @Description Create a new link with title, URL, and expiry date.
// @Tags links
// @Accept json
// @Produce json
// @Param link body entity.Link true "Link Data"
// @Success 201 {object} entity.Link
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /links [post]
func (h *LinkHandler) CreateLink(c *gin.Context) {
	var link entity.Link
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdLink, err := h.usecase.CreateLink(c.Request.Context(), &link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdLink)
}

// GetLink handles GET /links/:id
// GetLink godoc
// @Summary Get a link by ID
// @Description Retrieve a link using its ID.
// @Tags links
// @Accept json
// @Produce json
// @Param id path string true "Link ID"
// @Success 200 {object} entity.Link
// @Failure 404 {object} map[string]string "Link not found"
// @Security BearerAuth
// @Router /links/{id} [get]
func (h *LinkHandler) GetLink(c *gin.Context) {
	id := c.Param("id")
	link, err := h.usecase.GetLink(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}
	c.JSON(http.StatusOK, link)
}

// UpdateLink handles PUT /links/:id
// UpdateLink godoc
// @Summary Update an existing link
// @Description Update the link’s title, URL, or expiry date.
// @Tags links
// @Accept json
// @Produce json
// @Param id path string true "Link ID"
// @Param link body entity.Link true "Updated Link Data"
// @Success 200 {object} entity.Link
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /links/{id} [put]
func (h *LinkHandler) UpdateLink(c *gin.Context) {
	id := c.Param("id")
	var link entity.Link
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	link.ID = id // ensure ID matches the path param

	updatedLink, err := h.usecase.UpdateLink(c.Request.Context(), &link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedLink)
}

// DeleteLink handles DELETE /links/:id
// DeleteLink godoc
// @Summary Delete a link by ID
// @Description Delete the specified link.
// @Tags links
// @Accept json
// @Produce json
// @Param id path string true "Link ID"
// @Success 200 {object} map[string]string "Link deleted successfully"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /links/{id} [delete]
func (h *LinkHandler) DeleteLink(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteLink(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Link deleted successfully"})
}

// VisitLink handles GET /visit/:id
// VisitLink godoc
// @Summary Visit a link
// @Description Increment the link's click counter and log the visit.
// @Tags links
// @Accept json
// @Produce json
// @Param id path string true "Link ID"
// @Success 200 {object} entity.Link
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /visit/{id} [get]
func (h *LinkHandler) VisitLink(c *gin.Context) {
	id := c.Param("id")
	link, err := h.usecase.VisitLink(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, link)
}
