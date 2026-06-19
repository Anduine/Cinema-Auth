package domain

type UserResponse struct {
	Message string `json:"message"`
	Data    *User  `json:"data"`
}
