package response

type User struct {
	UUID     string  `json:"uuid"`
	Username string  `json:"username"`
	Name     string  `json:"name"`
	Groups   []Group `json:"groups"`
}
