package models

// Song represents a song in a concert setlist.
type Song struct {
    ID        int64  `json:"id"`
    Title     string `json:"title"`
    Notes     string `json:"notes"`
    ConcertID int64  `json:"concert_id"`
    Order     int    `json:"order"`
}
