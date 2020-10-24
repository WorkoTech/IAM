package requests;

import (
	"fmt"

	"worko.tech/iam/src/models"
	"worko.tech/iam/src/utils"

	"github.com/labstack/echo/v4"
)

type authRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}


type LoginRequest struct {
	authRequest
}
func (request *LoginRequest) Bind(ctx echo.Context) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	if err := ctx.Validate(request); err != nil {
		return err
	}
	return nil
}


type RegisterRequest struct {
	authRequest
}
func (request *RegisterRequest) Bind(ctx echo.Context, user *models.User) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	if err := ctx.Validate(request); err != nil {
		return err
	}

	fmt.Printf("Binding register request, password : %s\n", request.Password)
	user.Email = request.Email

	encryptedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}
	user.Password = encryptedPassword

	return nil
}


type SearchUsersRequest struct {
	UserIds []int64 `json:"userIds" validate:"required"`
}

func (request *SearchUsersRequest) Bind(ctx echo.Context) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	if err := ctx.Validate(request); err != nil {
		return err
	}
	return nil;
}
