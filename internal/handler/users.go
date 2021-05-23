package handler

import (
	"net/http"

	_ "github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/validator"
)

type usersSignUpInput struct {
	Username  string `json:"username" validator:"required,username,min=2,max=64"`
	FirstName string `json:"firstName" validator:"required,min=2,max=64"`
	LastName  string `json:"lastName" validator:"required,min=2,max=64"`
	Age       int    `json:"age" validator:"required,min=12,max=110"`
	Gender    int    `json:"gender" validator:"min=1,max=2"`
	Email     string `json:"email" validator:"required,email,max=64"`
	Password  string `json:"password" validator:"required,password,min=7,max=64"`
}

// @Summary Sign up
// @Tags users
// @ModuleID usersSignUp
// @Accept  json
// @Produce  json
// @Param input body usersSignUpInput true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404,409,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /users/sign-up [post]
func (h *Handler) usersSignUp(ctx *gorouter.Context) {
	var input usersSignUpInput

	if err := ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err := h.usersService.SignUp(service.UsersSignUpInput{
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Age:       input.Age,
		Gender:    input.Gender,
		Email:     input.Email,
		Password:  input.Password,
	})

	if err != nil {
		if err == service.ErrUserAlreadyExist {
			ctx.WriteError(http.StatusConflict, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusCreated)
}

type usersSignInInput struct {
	UsernameOrEmail string `json:"usernameOrEmail" validator:"required,max=64"`
	Password        string `json:"password" validator:"required,password,min=7,max=64"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// @Summary Sign in
// @Tags users
// @ModuleID usersSignIn
// @Accept  json
// @Produce  json
// @Param input body usersSignInInput true "sign in info"
// @Success 200 {object} tokenResponse
// @Failure 400,401,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /users/sign-in [post]
func (h *Handler) usersSignIn(ctx *gorouter.Context) {
	var input usersSignInInput

	if err := ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.usersService.SignIn(service.UsersSignInInput{
		UsernameOrEmail: input.UsernameOrEmail,
		Password:        input.Password,
	})

	if err != nil {
		if err == service.ErrUserWrongPassword {
			ctx.WriteError(http.StatusUnauthorized, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	ctx.WriteJSON(http.StatusOK, resp)
}

// @Summary Get user by ID
// @Security Auth
// @Tags users
// @ModuleID getUser
// @Accept  json
// @Produce  json
// @Param user_id path int true "ID of user"
// @Success 200 {object} model.User
// @Failure 400,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /users/{user_id} [GET]
func (h *Handler) getUser(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("user_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.usersService.GetByID(userID)
	if err != nil {
		if err == service.ErrUserDoesNotExist {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteJSON(http.StatusOK, user)
}

type usersRefreshTokensInput struct {
	AccessToken  string `json:"accessToken" validator:"required"`
	RefreshToken string `json:"refreshToken" validator:"required"`
}

// @Summary Refresh tokens
// @Tags users
// @ModuleID usersRefreshTokens
// @Accept  json
// @Produce  json
// @Param input body usersRefreshTokensInput true "tokens input"
// @Success 200 {object} tokenResponse
// @Failure 400,401,403,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /auth/refresh [POST]
func (h *Handler) usersRefreshTokens(ctx *gorouter.Context) {
	var input usersRefreshTokensInput

	if err := ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.usersService.RefreshTokens(service.UsersRefreshTokensInput{
		AccessToken:  input.AccessToken,
		RefreshToken: input.RefreshToken,
	})

	if err != nil {
		if err == service.ErrSessionNotFound {
			ctx.WriteError(http.StatusUnauthorized, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	ctx.WriteJSON(http.StatusOK, resp)
}

// @Summary Request moderator role
// @Security Auth
// @Tags users
// @ModuleID requestModerator
// @Accept  json
// @Produce  json
// @Success 201 {string} string "ok"
// @Failure 400,401,403,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /moderators/requests [POST]
func (h *Handler) requestModerator(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.usersService.CreateModeratorRequest(userID)
	if err != nil {
		if err == service.ErrModeratorRequestAlreadyExist {
			ctx.WriteError(http.StatusConflict, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusCreated)
}
