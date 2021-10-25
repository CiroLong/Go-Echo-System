package model

type User struct {
	ID       uint   `json:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
	Verified bool   `json:"verified"`
}
