package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"concerts/db"
	"concerts/models"
)

// ListSongs returns all songs for a specific concert.
func ListSongs(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    uid, ok := UserIDFromContext(ctx)
    if !ok {
        writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
        return
    }

    vars := mux.Vars(r)
    concertIDStr := vars["concertId"]
    concertID, err := strconv.ParseInt(concertIDStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, errors.New("invalid concert id"))
        return
    }

    // Verify the concert belongs to the user
    connection := db.Get()
    var concertUserID int64
    err = connection.QueryRow("SELECT user_id FROM concerts WHERE id = ?", concertID).Scan(&concertUserID)
    if err != nil {
        if err == sql.ErrNoRows {
            writeError(w, http.StatusNotFound, errors.New("concert not found"))
            return
        }
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db query error: %w", err))
        return
    }

    if concertUserID != uid {
        writeError(w, http.StatusForbidden, errors.New("access denied"))
        return
    }

    // Get songs for the concert
    rows, err := connection.Query("SELECT id, title, notes, concert_id, song_order FROM songs WHERE concert_id = ? ORDER BY song_order ASC, id ASC", concertID)
    if err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db query error: %w", err))
        return
    }
    defer rows.Close()

    var list []models.Song
    for rows.Next() {
        var song models.Song
        if err := rows.Scan(&song.ID, &song.Title, &song.Notes, &song.ConcertID, &song.Order); err != nil {
            writeError(w, http.StatusInternalServerError, fmt.Errorf("db scan error: %w", err))
            return
        }
        list = append(list, song)
    }

    writeJSON(w, http.StatusOK, list)
}

type createSongRequest struct {
    Title string `json:"title"`
    Notes string `json:"notes"`
}

// CreateSong inserts a new song for a specific concert.
func CreateSong(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    uid, ok := UserIDFromContext(ctx)
    if !ok {
        writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
        return
    }

    vars := mux.Vars(r)
    concertIDStr := vars["concertId"]
    concertID, err := strconv.ParseInt(concertIDStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, errors.New("invalid concert id"))
        return
    }

    // Verify the concert belongs to the user
    connection := db.Get()
    var concertUserID int64
    err = connection.QueryRow("SELECT user_id FROM concerts WHERE id = ?", concertID).Scan(&concertUserID)
    if err != nil {
        if err == sql.ErrNoRows {
            writeError(w, http.StatusNotFound, errors.New("concert not found"))
            return
        }
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db query error: %w", err))
        return
    }

    if concertUserID != uid {
        writeError(w, http.StatusForbidden, errors.New("access denied"))
        return
    }

    var req createSongRequest
    if err := readJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
        return
    }

    if req.Title == "" {
        writeError(w, http.StatusBadRequest, errors.New("title is required"))
        return
    }

    // Get the next order number
    var maxOrder int
    err = connection.QueryRow("SELECT COALESCE(MAX(song_order), -1) FROM songs WHERE concert_id = ?", concertID).Scan(&maxOrder)
    if err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db query error: %w", err))
        return
    }

    res, err := connection.Exec("INSERT INTO songs (title, notes, concert_id, song_order) VALUES (?, ?, ?, ?)", 
        req.Title, req.Notes, concertID, maxOrder+1)
    if err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db insert error: %w", err))
        return
    }

    id, _ := res.LastInsertId()
    writeJSON(w, http.StatusCreated, models.Song{
        ID:        id,
        Title:     req.Title,
        Notes:     req.Notes,
        ConcertID: concertID,
        Order:     maxOrder + 1,
    })
}

// DeleteSong deletes a song by id.
func DeleteSong(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    uid, ok := UserIDFromContext(ctx)
    if !ok {
        writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
        return
    }

    vars := mux.Vars(r)
    songIDStr := vars["songId"]
    songID, err := strconv.ParseInt(songIDStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, errors.New("invalid song id"))
        return
    }

    connection := db.Get()
    
    // Delete the song and verify it belonged to the user
    res, err := connection.Exec(`
        DELETE FROM songs 
        WHERE id = ? AND concert_id IN (
            SELECT id FROM concerts WHERE user_id = ?
        )`, songID, uid)
    if err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db delete error: %w", err))
        return
    }

    n, _ := res.RowsAffected()
    if n == 0 {
        writeError(w, http.StatusNotFound, sql.ErrNoRows)
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{"deleted": songID})
}

// UpdateSongOrder updates the order of songs in a setlist.
func UpdateSongOrder(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    uid, ok := UserIDFromContext(ctx)
    if !ok {
        writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
        return
    }

    vars := mux.Vars(r)
    concertIDStr := vars["concertId"]
    concertID, err := strconv.ParseInt(concertIDStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, errors.New("invalid concert id"))
        return
    }

    // Verify the concert belongs to the user
    connection := db.Get()
    var concertUserID int64
    err = connection.QueryRow("SELECT user_id FROM concerts WHERE id = ?", concertID).Scan(&concertUserID)
    if err != nil {
        if err == sql.ErrNoRows {
            writeError(w, http.StatusNotFound, errors.New("concert not found"))
            return
        }
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db query error: %w", err))
        return
    }

    if concertUserID != uid {
        writeError(w, http.StatusForbidden, errors.New("access denied"))
        return
    }

    type songOrderUpdate struct {
        SongID int64 `json:"song_id"`
        Order  int   `json:"order"`
    }

    var updates []songOrderUpdate
    if err := readJSON(r, &updates); err != nil {
        writeError(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
        return
    }

    // Start a transaction
    tx, err := connection.Begin()
    if err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db transaction error: %w", err))
        return
    }
    defer tx.Rollback()

    // Update each song's order
    for _, update := range updates {
        _, err := tx.Exec("UPDATE songs SET song_order = ? WHERE id = ? AND concert_id = ?", 
            update.Order, update.SongID, concertID)
        if err != nil {
            writeError(w, http.StatusInternalServerError, fmt.Errorf("db update error: %w", err))
            return
        }
    }

    if err := tx.Commit(); err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db commit error: %w", err))
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{"updated": len(updates)})
}
