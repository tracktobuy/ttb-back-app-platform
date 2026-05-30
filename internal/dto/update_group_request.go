package dto

type UpdateGroupRequest struct {
	Name           string  `json:"name"`
	Budget         float32 `json:"budget"`
	BudgetCurrency string  `json:"budgetCurrency"`
}
