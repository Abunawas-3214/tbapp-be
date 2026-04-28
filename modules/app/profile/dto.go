package profile

type StoreProfileResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Address     *string `json:"address"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	TaxID       *string `json:"tax_id"`
	LogoURL     *string `json:"logo_url"`
}
