package sysadmin

type CreateAdminRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Level    string `json:"level" validate:"required"`
}

type AdminResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Level string `json:"level"`
}

type AdminListResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Level      string `json:"level"`
	IsActive   bool   `json:"is_active"`
	AdminSince string `json:"admin_since"`
}

type UpdateAdminRequest struct {
	Name     *string `json:"name"`
	Level    *string `json:"level"`
	IsActive *bool   `json:"is_active"`
}

type ChangePasswordRequest struct {
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type ChangeEmailRequest struct {
	NewEmail string `json:"new_email" validate:"required,email"`
}
