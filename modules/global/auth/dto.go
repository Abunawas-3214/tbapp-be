package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type StoreAccessDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	SchemaName string `json:"schema_name"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID         string  `json:"id"`
		Name       string  `json:"name"`
		Email      string  `json:"email"`
		AdminLevel *string `json:"admin_level,omitempty"`
	} `json:"user"`
	Stores []StoreAccessDTO `json:"stores,omitempty"`
}

type SelectStoreRequest struct {
	StoreID string `json:"store_id" validate:"required"`
}

type SelectStoreResponse struct {
	Token      string `json:"token"`
	SchemaName string `json:"schema_name"`
	Message    string `json:"message"`
}
