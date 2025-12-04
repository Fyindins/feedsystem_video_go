package api

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
}

type RenameByIDRequest struct {
	ID          uint   `json:"id"`
	NewUsername string `json:"new_username"`
}

type RenameByIDResponse struct {
}

type FindByIDRequest struct {
	ID uint `json:"id"`
}

type FindByIDResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type FindByUsernameRequest struct {
	Username string `json:"username"`
}

type FindByUsernameResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type ChangePasswordRequest struct {
	ID          uint   `json:"id"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct {
}
