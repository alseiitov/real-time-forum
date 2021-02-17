package handler

import (
	"net/http"

	"github.com/alseiitov/real-time-forum/internal/service"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/validator"
)

type usersSignUpInput struct {
	Username  string `json:"username" validator:"required,username,min=2,max=64"`
	FirstName string `json:"firstName" validator:"required,min=2,max=64"`
	LastName  string `json:"lastName" validator:"required,min=2,max=64"`
	Age       int    `json:"age" validator:"required,min=12,max=110"`
	Gender    int    `json:"gender" validator:"min=0,max=1"`
	Email     string `json:"email" validator:"required,email,max=64"`
	Password  string `json:"password" validator:"required,password,min=7,max=64"`
}

func (h *Handler) usersSignUp(ctx *gorouter.Context) {
	var input usersSignUpInput

	err := ctx.ReadBody(&input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = validator.Validate(input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.usersService.SignUp(service.UsersSignUpInput{
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Age:       input.Age,
		Gender:    input.Gender,
		Email:     input.Email,
		Password:  input.Password,
	})

	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	ctx.ResponseWriter.WriteHeader(http.StatusCreated)
}

type usersSignInInput struct {
	UsernameOrEmail string `json:"usernameOrEmail" validator:"required,max=64"`
	Password        string `json:"password" validator:"required,password,min=7,max=64"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) usersSignIn(ctx *gorouter.Context) {
	var input usersSignInInput

	err := ctx.ReadBody(&input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = validator.Validate(input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.usersService.SignIn(service.UsersSignInInput{
		UsernameOrEmail: input.UsernameOrEmail,
		Password:        input.Password,
	})
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *Handler) getUser(ctx *gorouter.Context) {

}

func (h *Handler) updateUser(ctx *gorouter.Context) {

}

type usersRefreshTokensInput struct {
	AccessToken  string `json:"accessToken" validator:"required"`
	RefreshToken string `json:"refreshToken" validator:"required"`
}

func (h *Handler) usersRefreshTokens(ctx *gorouter.Context) {
	var input usersRefreshTokensInput

	err := ctx.ReadBody(&input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = validator.Validate(input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.usersService.RefreshTokens(service.UsersRefreshTokensInput{
		AccessToken:  input.AccessToken,
		RefreshToken: input.RefreshToken,
	})
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
