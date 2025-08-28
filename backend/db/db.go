package db

import (
    "database/sql"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "sync"

    _ "modernc.org/sqlite"
)

var (
    conn *sql.DB
    once sync.Once
)

// Init initializes the SQLite database connection and runs migrations.
func Init() (*sql.DB, error) {
    var initErr error
    once.Do(func() {
        dbPath := getDBPath()
        if err := ensureDir(filepath.Dir(dbPath)); err != nil {
            initErr = fmt.Errorf("failed to ensure db directory: %w", err)
            return
        }

        // modernc.org/sqlite registers the driver as "sqlite"
        dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbPath)
        c, err := sql.Open("sqlite", dsn)
        if err != nil {
            initErr = fmt.Errorf("failed to open sqlite: %w", err)
            return
        }

        // Reasonable defaults for modernc sqlite
        if _, err := c.Exec("PRAGMA foreign_keys = ON;"); err != nil {
            _ = c.Close()
            initErr = fmt.Errorf("failed to enable foreign_keys: %w", err)
            return
        }

        if err := migrate(c); err != nil {
            _ = c.Close()
            initErr = fmt.Errorf("failed to run migrations: %w", err)
            return
        }

        conn = c
    })
    return conn, initErr
}

// Get returns the initialized sql.DB. Panics if Init was not called.
func Get() *sql.DB {
    if conn == nil {
        panic("db not initialized: call db.Init() first")
    }
    return conn
}

func getDBPath() string {
    if p := os.Getenv("DB_PATH"); p != "" {
        return p
    }
    return filepath.Join("data", "concerts.db")
}

func ensureDir(dir string) error {
    if dir == "." || dir == "" {
        return nil
    }
    if st, err := os.Stat(dir); err == nil {
        if !st.IsDir() {
            return errors.New("db directory path is not a directory")
        }
        return nil
    }
    return os.MkdirAll(dir, 0o755)
}

func migrate(c *sql.DB) error {
    stmts := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            password_hash TEXT NOT NULL
        );`,
        `CREATE TABLE IF NOT EXISTS concerts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            date TEXT NOT NULL,
            location TEXT NOT NULL,
            user_id INTEGER NOT NULL,
            FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
        );`,
        `CREATE INDEX IF NOT EXISTS idx_concerts_user_id ON concerts(user_id);`,
    }
    for _, s := range stmts {
        if _, err := c.Exec(s); err != nil {
            return fmt.Errorf("migration failed: %w", err)
        }
    }
    return nil
}


