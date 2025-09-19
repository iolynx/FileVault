package util

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

// NewText converts a string to a pgtype.Text.
func NewText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

// TimeToTimestamptz converts a *time.Time to a pgtype.Timestamptz.
func TimeToTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// BoolToPgBool converts a *bool to a pgtype.Bool.
func BoolToPgBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}
