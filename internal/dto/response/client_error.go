package response

type ClientError struct {
	Status    int                `json:"status"`
	Code      string             `json:"code"`
	Message   string             `json:"message"`
	Details   ClientErrorDetails `json:"details"`
	Timestamp string             `json:"timestamp"`
	Path      string             `json:"path"`
}

type ClientErrorDetails struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}
