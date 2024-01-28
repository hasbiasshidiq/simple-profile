// This file contains types that are used in the repository layer.
package repository

import "time"

// Profile, representing profile object on repository
type Profile struct {
	ID          uint64
	FullName    string `json:"full_name" `
	CountryCode string `json:"country_code" `
	PhoneNumber string `json:"phone_number" `
	Password    string `json:"password"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
