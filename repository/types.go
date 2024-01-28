// This file contains types that are used in the repository layer.
package repository

import "time"

// Profile, representing profile object on repository
type Profile struct {
	ID          uint64     `json:"id" db:"id"`
	FullName    string     `json:"full_name" db:"full_name"`
	CountryCode string     `json:"country_code"`
	PhoneNumber string     `json:"phone_number" db:"phone_number"`
	Password    string     `json:"password"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
