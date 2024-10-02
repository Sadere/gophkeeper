// Provides models used both by server and client
package model

import (
	"time"
)

type Secret struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	UserID    uint64    `db:"user_id"`
	Metadata  string    `db:"metadata"`
	SType     string    `db:"ent_type"`
	Payload   []byte    `db:"payload"`

	Creds *Credentials `db:"-"`
	Text  *Text        `db:"-"`
	Blob  *Blob        `db:"-"`
	Card  *Card        `db:"-"`
}

type Secrets []*Secret

// Secret type
type SecretType string

const (
	CredSecret    SecretType = "credential"
	TextSecret    SecretType = "text"
	BlobSecret    SecretType = "blob"
	CardSecret    SecretType = "card"
	UnknownSecret SecretType = "unknown"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Text struct {
	Content string `json:"content"`
}

type Blob struct {
	FileName string `json:"file_name"`
	IsDone   bool   `json:"is_done"`
}

type Card struct {
	Number   string `json:"number"`
	ExpYear  uint32 `json:"exp_year"`
	ExpMonth uint32 `json:"exp_month"`
	Cvv      uint32 `json:"cvv"`
}

// Indicates status of secret for particular client: new, updated or read
type SecretPreviewStatus string

const (
	SecretPreviewNew     SecretPreviewStatus = "new"
	SecretPreviewUpdated SecretPreviewStatus = "updated"
	SecretPreviewRead    SecretPreviewStatus = "read"
)

// Holds just preview secret information: metadata, date, id.
//
// Doesn't include any private user info
type SecretPreview struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	Metadata  string
	SType     string
	Status    SecretPreviewStatus
}

type SecretPreviews []*SecretPreview
