package handler

import (
	"net/http"
	"strconv"

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
	categoryIDParam, _ := ctx.GetParam("category_id")
	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	pageParam, _ := ctx.GetParam("page")
	page, err := strconv.Atoi(pageParam)
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
