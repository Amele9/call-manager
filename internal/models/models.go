package models

import "time"

// CallInfo is the call information
type CallInfo struct {
	// ID is the call ID from the database
	ID int `json:"id" binding:"-"`

	// ClientName is the name of the client who created the call
	ClientName string `json:"client_name" binding:"required"`

	// PhoneNumber is the phone number of the client who created the call
	PhoneNumber string `json:"phone_number" binding:"required,e164"`

	// Description is a description of the call
	Description string `json:"description" binding:"required"`

	// Status is the status of the call
	// The status takes on 2 values: open, closed
	Status string `json:"status" binding:"oneof=open closed"`

	// CreatedAt is the time the call was created
	CreatedAt time.Time `json:"created_at" binding:"-"`
}
