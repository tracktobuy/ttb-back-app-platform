package domain

import "time"

type User struct {
	UUID      string    `json:"uuid"`
	Version   int       `json:"-"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
