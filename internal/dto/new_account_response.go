package dto

import "github.com/tracktobuy/ttb-back-app-platform/internal/domain"

type userResponse struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type groupResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type newAccountResponse struct {
	User  userResponse  `json:"user"`
	Group groupResponse `json:"group"`
}

func NewAccountResponse(user domain.User, group domain.Group) newAccountResponse {
	return newAccountResponse{
		User: userResponse{
			UUID:     user.UUID,
			Username: user.Username,
			Name:     user.Name,
		},
		Group: groupResponse{
			UUID: group.UUID,
			Name: group.Name,
		},
	}
}
