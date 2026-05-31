package dto

type GroupResponse struct {
	UUID           string  `json:"uuid"`
	Name           string  `json:"name"`
	Budget         float32 `json:"budget"`
	BudgetCurrency string  `json:"budgetCurrency"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
	CreatedBy      string  `json:"createdBy,omitempty"`
}
