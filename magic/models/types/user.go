package types

type LoginRequest struct {
	Email  string `json:"email"`
	Passwd string `json:"passwd"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
	Code     string `json:"code"`
}

type SendEmailCodeRequest struct {
	Email string `json:"email"`
}

type DeleteUserRequest struct {
	Uid int `json:"uid"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UpdateProfileRequest struct {
	Gender   string `json:"gender,omitempty"`
	Birthday string `json:"birthday,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Name     string `json:"name,omitempty"`
}

type AdminUpdateUserRequest struct {
	Uid      int      `json:"uid"`
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	IsActive *int     `json:"isActive,omitempty"`
}
