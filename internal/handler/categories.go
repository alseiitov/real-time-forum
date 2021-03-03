package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
)

func (h *Handler) getAllCategories(ctx *gorouter.Context) {
	categories, err := h.categoriesService.GetAll()
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, categories)
}

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
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	category.Posts, err = h.postsService.GetPostsByCategoryID(categoryID, page)
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, category)
}
