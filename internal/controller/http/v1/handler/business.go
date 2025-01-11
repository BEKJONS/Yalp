package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"yalp_ulab/config"
	"yalp_ulab/internal/entity"
)

// CreateBusiness godoc
// @Router /business [post]
// @Summary Create a new business
// @Description Create a new business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param business body entity.Business true "Business object"
// @Success 201 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateBusiness(ctx *gin.Context) {
	var body entity.Business

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", http.StatusBadRequest)
		return
	}

	business, err := h.UseCase.BusinessRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating business") {
		return
	}

	ctx.JSON(http.StatusCreated, business)
}

// GetBusiness godoc
// @Router /business/{id} [get]
// @Summary Get a business by ID
// @Description Get a business by ID
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Success 200 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBusiness(ctx *gin.Context) {
	var req entity.BusinessSingleRequest

	req.ID = ctx.Param("id")

	business, err := h.UseCase.BusinessRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting business") {
		return
	}

	ctx.JSON(http.StatusOK, business)
}

// GetBusinesses godoc
// @Router /business/list [get]
// @Summary Get a list of businesses
// @Description Get a list of businesses
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.BusinessList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBusinesses(ctx *gin.Context) {
	var req entity.GetListFilter

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "business_name",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "category",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "location",
			Type:   "search",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	businesses, err := h.UseCase.BusinessRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting businesses") {
		return
	}

	ctx.JSON(http.StatusOK, businesses)
}

// UpdateBusiness godoc
// @Router /business [put]
// @Summary Update a business
// @Description Update a business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param business body entity.Business true "Business object"
// @Success 200 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateBusiness(ctx *gin.Context) {
	var body entity.Business

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", http.StatusBadRequest)
		return
	}

	business, err := h.UseCase.BusinessRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating business") {
		return
	}

	ctx.JSON(http.StatusOK, business)
}

// DeleteBusiness godoc
// @Router /business/{id} [delete]
// @Summary Delete a business
// @Description Delete a business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteBusiness(ctx *gin.Context) {
	var req entity.Id

	req.ID = ctx.Param("id")

	err := h.UseCase.BusinessRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting business") {
		return
	}

	ctx.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Business deleted successfully",
	})
}
