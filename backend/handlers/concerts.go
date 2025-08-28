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

// ListConcerts returns all concerts for the authenticated user.
func ListConcerts(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    uid, ok := UserIDFromContext(ctx)
    if !ok {
        writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
        return
    }
    connection := db.Get()
    rows, err := connection.Query("SELECT id, title, date, location, user_id FROM concerts WHERE user_id = ? ORDER BY date DESC", uid)
    if err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db query error: %w", err))
        return
    }
    defer rows.Close()
    var list []models.Concert
    for rows.Next() {
        var c models.Concert
        if err := rows.Scan(&c.ID, &c.Title, &c.Date, &c.Location, &c.UserID); err != nil {
            writeError(w, http.StatusInternalServerError, fmt.Errorf("db scan error: %w", err))
            return
        }
        list = append(list, c)
    }
    writeJSON(w, http.StatusOK, list)
}

func GetConcert(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		uid, ok := UserIDFromContext(ctx)
		if !ok {
				writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
				return
		}
		vars := mux.Vars(r)
		idStr := vars["id"]
		cid, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
				writeError(w, http.StatusBadRequest, errors.New("invalid id"))
				return
		}
		connection := db.Get()
		var c models.Concert
		if err := connection.QueryRow("SELECT id, title, date, location, user_id FROM concerts WHERE id = ? AND user_id = ?", cid, uid).Scan(&c.ID, &c.Title, &c.Date, &c.Location, &c.UserID); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Errorf("db query error: %w", err))
				return
		}
		writeJSON(w, http.StatusOK, c)
}


type createConcertRequest struct {
    Title    string `json:"title"`
    Date     string `json:"date"`
    Location string `json:"location"`
}

// CreateConcert inserts a new concert for the authenticated user.
func CreateConcert(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    uid, ok := UserIDFromContext(ctx)
    if !ok {
        writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
        return
    }
    var req createConcertRequest
    if err := readJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
        return
    }
    if req.Title == "" || req.Date == "" || req.Location == "" {
        writeError(w, http.StatusBadRequest, errors.New("title, date, and location are required"))
        return
    }
    connection := db.Get()
    res, err := connection.Exec("INSERT INTO concerts (title, date, location, user_id) VALUES (?, ?, ?, ?)", req.Title, req.Date, req.Location, uid)
    if err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db insert error: %w", err))
        return
    }
    id, _ := res.LastInsertId()
    writeJSON(w, http.StatusCreated, models.Concert{
        ID:       id,
        Title:    req.Title,
        Date:     req.Date,
        Location: req.Location,
        UserID:   uid,
    })
}

// DeleteConcert deletes a concert by id for the authenticated user.
func DeleteConcert(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    uid, ok := UserIDFromContext(ctx)
    if !ok {
        writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
        return
    }
    vars := mux.Vars(r)
    idStr := vars["id"]
    cid, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, errors.New("invalid id"))
        return
    }
    connection := db.Get()
    res, err := connection.Exec("DELETE FROM concerts WHERE id = ? AND user_id = ?", cid, uid)
    if err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db delete error: %w", err))
        return
    }
    n, _ := res.RowsAffected()
    if n == 0 {
        writeError(w, http.StatusNotFound, sql.ErrNoRows)
        return
    }
    writeJSON(w, http.StatusOK, map[string]any{"deleted": cid})
}


