package response

type Item struct {
	UUID      string   `json:"uuid"`
	Title     string   `json:"title"`
	Images    []string `json:"images,omitempty"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt,omitempty"`
	CreatedBy string   `json:"createdBy"`
	Labels    []string `json:"labels,omitempty"`
	Stores    []string `json:"stores,omitempty"`
	Groups    []string `json:"groups,omitempty"`
}
