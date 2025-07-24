package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role   string `json:"role"` // e.g., "admin", "user"
}
