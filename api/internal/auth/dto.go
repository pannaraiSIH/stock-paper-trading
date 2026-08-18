package auth

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type CreateUserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserResponse struct {
	AccessToken string `json:"accessToken"`
}
