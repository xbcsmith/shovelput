package main

import (
	"crypto/rand"
	"os"
	"time"

	"github.com/oklog/ulid"
)

// GetEnv returns an env variable value or a default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// NewULID returns a ULID as a string.
func NewULID() string {
	newid, _ := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	return newid.String()
}
