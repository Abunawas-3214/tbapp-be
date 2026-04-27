package tenant

type CreateTenantRequest struct {
	// --- Data untuk skema PUBLIC (Data Pusat) ---
	StoreName  string `json:"store_name" validate:"required"`
	StoreSlug  string `json:"store_slug" validate:"required,alphanum"`
	SchemaName string `json:"schema_name" validate:"required"`

	// --- Data untuk skema TENANT & Identitas Owner ---
	OwnerName     string `json:"owner_name" validate:"required"`
	OwnerEmail    string `json:"owner_email" validate:"required,email"`
	OwnerPassword string `json:"owner_password" validate:"required,min=8"`
}

type CreateTenantResponse struct {
	StoreID   string `json:"store_id"`
	StoreName string `json:"store_name"`
	StoreSlug string `json:"store_slug"`
	OwnerID   string `json:"owner_id"`
	OwnerName string `json:"owner_name"`
}
