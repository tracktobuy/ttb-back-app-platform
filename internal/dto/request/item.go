package request

type ItemRequest struct {
	URL       string   `json:"url"`
	Title     string   `json:"title"`
	Price     float32  `json:"price"`
	Currency  string   `json:"currency"`
	Image     string   `json:"image"`
	Domain    string   `json:"domain"`
	Store     string   `json:"store"`
	Labels    []string `json:"labels"`
	GroupId   string   `json:"groupId"`
	CreatedBy string   `json:"createdBy"`
}
