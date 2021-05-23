package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
	_ "github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"
)

// @Summary Get list of all categories
// @Security Auth
// @Tags categories
// @ModuleID getAllCategories
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Category "ok"
// @Failure 400,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /categories [GET]
func (h *Handler) getAllCategories(ctx *gorouter.Context) {
	categories, err := h.categoriesService.GetAll()
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, categories)
}

// @Summary Get list of N posts of category page
// @Security Auth
// @Tags categories
// @ModuleID getCategoryPage
// @Accept  json
// @Produce  json
// @Param category_id path int true "ID of category"
// @Param page path int true "page number"
// @Success 200 {object} model.Category
// @Failure 400,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /categories/{category_id}/{page} [GET]
func (h *Handler) getCategoryPage(ctx *gorouter.Context) {
	categoryID, err := ctx.GetIntParam("category_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	page, err := ctx.GetIntParam("page")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.categoriesService.GetByID(categoryID, page)
	if err != nil {
		if err == service.ErrCategoryDoesntExist {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteJSON(http.StatusOK, category)
}
