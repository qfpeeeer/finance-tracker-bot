package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// UserState represents information about a user state entry.
type UserState struct {
	db *sqlx.DB
}

// UserStateInfo represents the structure of a user's state information.
type UserStateInfo struct {
	ID        int64                  `db:"id"`
	UserID    int64                  `db:"user_id"`
	State     string                 `db:"state"`
	DataJSON  string                 `db:"data"` // Store as JSON
	Data      map[string]interface{} `db:"-"`    // Don't store in DB, use for application logic
	Timestamp time.Time              `db:"timestamp"`
}

// NewUserState creates a new UserState storage
func NewUserState(db *sqlx.DB) (*UserState, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user_states (
		id INTEGER PRIMARY KEY,
		user_id INTEGER UNIQUE,
		state TEXT,
		data TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create user_states table: %w", err)
	}

	// Add index on user_id for faster lookup
	if _, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_states_user_id ON user_states(user_id)`); err != nil {
		return nil, fmt.Errorf("failed to create index on user_id: %w", err)
	}

	return &UserState{db: db}, nil
}

// Write adds or updates a user's state entry
func (us *UserState) Write(entry UserStateInfo) error {
	query := `INSERT INTO user_states (user_id, state, data) VALUES (?, ?, ?) ON CONFLICT(user_id) DO UPDATE SET state = excluded.state, data = excluded.data`
	if _, err := us.db.Exec(query, entry.UserID, entry.State, entry.DataJSON); err != nil {
		return fmt.Errorf("failed to insert or update user state entry: %w", err)
	}

	log.Printf("[info] User state updated for user_id: %d, state: %s", entry.UserID, entry.State)
	return nil
}

// Read returns the latest state entry for a given user ID
func (us *UserState) Read(userID int64) (*UserStateInfo, error) {
	var entry UserStateInfo
	err := us.db.Get(&entry, "SELECT * FROM user_states WHERE user_id = ? ORDER BY timestamp DESC LIMIT 1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user state entry: %w", err)
	}

	entry.Timestamp = entry.Timestamp.Local()
	return &entry, nil
}
