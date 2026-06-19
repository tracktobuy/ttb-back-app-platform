package response

type ServerError struct {
	Status    int                `json:"status"`
	Code      string             `json:"code"`
	Message   string             `json:"message"`
	Details   ServerErrorDetails `json:"details"`
	Timestamp string             `json:"timestamp"`
	Path      string             `json:"path"`
}

type ServerErrorDetails struct {
	Service    string `json:"service"`
	RetryAfter int    `json:"retryAfter"`
}
