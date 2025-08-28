package handlers

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"

    "concerts/db"
)

type contextKey string

const userIDContextKey contextKey = "uid"

func getJWTSecret() []byte {
    if s := os.Getenv("JWT_SECRET"); s != "" {
        return []byte(s)
    }
    return []byte("dev-secret-change-me")
}

type registerRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type loginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Register creates a new user with a hashed password.
func Register(w http.ResponseWriter, r *http.Request) {
    var req registerRequest
    if err := readJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
        return
    }
    if len(strings.TrimSpace(req.Username)) < 3 || len(req.Password) < 6 {
        writeError(w, http.StatusBadRequest, errors.New("username must be >=3 and password >=6 characters"))
        return
    }

    connection := db.Get()
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to hash password: %w", err))
        return
    }
    _, err = connection.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", req.Username, string(hashed))
    if err != nil {
        // crude unique detection
        if strings.Contains(strings.ToLower(err.Error()), "unique") {
            writeError(w, http.StatusConflict, errors.New("username already exists"))
            return
        }
        writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to insert user: %w", err))
        return
    }
    writeJSON(w, http.StatusCreated, map[string]string{"message": "registered"})
}

// Login verifies credentials and returns a JWT token.
func Login(w http.ResponseWriter, r *http.Request) {
    var req loginRequest
    if err := readJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
        return
    }
    connection := db.Get()
    var (
        id int64
        hash string
    )
    err := connection.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", req.Username).Scan(&id, &hash)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            writeError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
            return
        }
        writeError(w, http.StatusInternalServerError, fmt.Errorf("db error: %w", err))
        return
    }
    if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
        writeError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
        return
    }

    claims := jwt.RegisteredClaims{
        Subject:   strconv.FormatInt(id, 10),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString(getJWTSecret())
    if err != nil {
        writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to sign token: %w", err))
        return
    }
    writeJSON(w, http.StatusOK, map[string]any{
        "token": signed,
        "user": map[string]any{
            "id":       id,
            "username": req.Username,
        },
    })
}

// RequireAuth validates JWT from Authorization header and injects user id into context.
func RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        header := r.Header.Get("Authorization")
        if header == "" || !strings.HasPrefix(header, "Bearer ") {
            writeError(w, http.StatusUnauthorized, errors.New("missing or invalid authorization header"))
            return
        }
        tokenString := strings.TrimPrefix(header, "Bearer ")
        parsed, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
            if t.Method != jwt.SigningMethodHS256 {
                return nil, errors.New("unexpected signing method")
            }
            return getJWTSecret(), nil
        })
        if err != nil || !parsed.Valid {
            writeError(w, http.StatusUnauthorized, errors.New("invalid token"))
            return
        }
        claims, ok := parsed.Claims.(*jwt.RegisteredClaims)
        if !ok || claims.Subject == "" {
            writeError(w, http.StatusUnauthorized, errors.New("invalid token claims"))
            return
        }
        uid, err := strconv.ParseInt(claims.Subject, 10, 64)
        if err != nil {
            writeError(w, http.StatusUnauthorized, errors.New("invalid subject"))
            return
        }
        ctx := context.WithValue(r.Context(), userIDContextKey, uid)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// thin wrappers moved to auth_ctx.go


