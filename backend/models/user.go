package models

// User represents an application user.
type User struct {
    ID           int64  `json:"id"`
    Username     string `json:"username"`
    PasswordHash string `json:"-"`
}


