package delivery

// Types for request and responses
type (
	// CreateRegisterRequest struct
	CreateUserRequest struct {
		Email     string `json:"email"`
		Passwords string `json:"passwords"`
	}
	// CreateRegisterResponse struct
	CreateUserResponse struct {
		Status string `json:"status"`
	}
	// CreateLoginRequest struct
	GetUserRequest struct {
		Id string `json:"id"`
	}
	// CreateLoginResponse struct
	GetUserResponse struct {
		Status string `json:"status"`
	}
)
