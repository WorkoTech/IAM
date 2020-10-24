package responses;

import (
	"worko.tech/iam/src/models"
)

type userResponse struct {
	ID 			uint `json:"id"`
	Email 		string `json:"email"`
}

func NewUserResponse(user *models.User) *userResponse {
	response := new(userResponse)
	response.ID = user.ID
	response.Email = user.Email

	return response
}

func NewUsersResponse(users *[]models.User) []*userResponse {
	response := make([]*userResponse, len(*users))

	for i, user := range *users {
		response[i] = new(userResponse)
		response[i].ID = user.ID
		response[i].Email = user.Email
	}

	return response
}

type accessUserResponse struct {
	ID 			uint `json:"id"`
	Email 		string `json:"email"`
	AccessToken string `json:"accessToken"`
}

func NewAccessUserResponse(user *models.User, token string) *accessUserResponse {
	response := new(accessUserResponse)
	response.ID = user.ID
	response.Email = user.Email
	response.AccessToken = token

	return response
}
