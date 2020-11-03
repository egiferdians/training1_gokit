package delivery

// Types for request and responses
type (
	// CreateUserRequest struct
	CreateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// CreateRegisterResponse struct
	CreateUserResponse struct {
		Ok string `json:"ok"`
	}
	// GetUserRequest struct
	GetUserRequest struct {
		Id string `json:"id"`
	}
	// CreateLoginResponse struct
	GetUserResponse struct {
		Email string `json:"email"`
	}
)
