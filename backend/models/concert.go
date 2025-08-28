package models

// Concert represents a concert record owned by a user.
type Concert struct {
    ID       int64  `json:"id"`
    Title    string `json:"title"`
    Date     string `json:"date"`
    Location string `json:"location"`
    UserID   int64  `json:"user_id"`
}


