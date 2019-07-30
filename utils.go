package main

import (
	"crypto/rand"
	"github.com/oklog/ulid"
	"time"
)

// NewULID returns a ULID as a string.
func NewULID() string {
	newid, _ := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	return newid.String()
}
