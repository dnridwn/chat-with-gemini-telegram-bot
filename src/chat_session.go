package main

import (
	"database/sql"
	"encoding/json"
)

type ChatSession struct {
	ID      uint64 `json:"id,omitempty"`
	History string `json:"history,omitempty"`
	// CreatedAt sql.NullTime
	// UpdatedAt sql.NullTime
	// DeletedAt sql.NullTime
}

func GetLatestHistory(db *sql.DB, cID int64) (ChatSession, error) {
	row := db.QueryRow("SELECT id, history FROM chat_sessions WHERE deleted_at IS NULL AND chat_id = ? ORDER BY ID DESC LIMIT 1", cID)
	if row.Err() != nil {
		if row.Err() != sql.ErrNoRows {
			return ChatSession{}, nil
		}

		return ChatSession{}, row.Err()
	}

	c := ChatSession{}
	if err := row.Scan(&c.ID, &c.History); err != nil {
		if row.Err() != sql.ErrNoRows {
			return ChatSession{}, nil
		}

		return ChatSession{}, err
	}

	return c, nil
}

func SaveHistory(db *sql.DB, cID int64, cs []Content) error {
	csB, err := json.Marshal(cs)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO chat_sessions (chat_id, history, created_at, updated_at) VALUES (?, ?, NOW(), NOW())", cID, string(csB))
	if err != nil {
		return err
	}

	return nil
}
