package models

import "github.com/google/uuid"

type Good struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	OwnerID   *uuid.UUID
	Name      string
}

type Bad struct {
	ID        string    // want `field ID must be uuid\.UUID`
	AccountID string    // want `field AccountID must be uuid\.UUID`
	Name      string
}
