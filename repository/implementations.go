package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func (r *Repository) CreateProfile(input Profile) (createdID int, err error) {

	stmt, err := r.Db.Prepare(
		`INSERT INTO profiles
			(
				full_name, 
				country_code, 
				phone_number, 
				password
			) VALUES ($1, $2, $3, $4) RETURNING id`)

	if err != nil {
		return createdID, err
	}

	err = stmt.QueryRow(
		input.FullName,
		input.CountryCode,
		input.PhoneNumber,
		input.Password,
	).Scan(&createdID)

	if err != nil {
		return createdID, err
	}

	return createdID, nil
}

func (r *Repository) GetPhoneNumberExistence(phoneNumber string) (isExist bool, err error) {
	var profileID int
	err = r.Db.QueryRow(`
		SELECT 
			id 
		FROM 
			profiles
		WHERE 
			phone_number = $1 and deleted_at is null`,
		phoneNumber).Scan(&profileID)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err

}

func (r *Repository) GetPhoneNumberExistenceWithExcludedID(phoneNumber string, excludedID int) (isExist bool, err error) {
	var profileID int
	err = r.Db.QueryRow(`
		SELECT 
			id 
		FROM 
			profiles
		WHERE 
			phone_number = $1 and id != $2 and deleted_at is null`,
		phoneNumber, excludedID).Scan(&profileID)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err

}
func (r *Repository) GetProfileByPhoneNumber(phoneNumber string) (profile Profile, err error) {

	// Fetch a single row from the database
	row := r.Db.QueryRow(`
		SELECT * FROM 
			profiles 
		WHERE 
			phone_number = $1 and deleted_at is null`, phoneNumber)

	err = row.Scan(
		&profile.ID,
		&profile.FullName,
		&profile.CountryCode,
		&profile.PhoneNumber,
		&profile.Password,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.DeletedAt,
	)
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (r *Repository) GetProfileByID(id int) (profile Profile, err error) {

	// Fetch a single row from the database
	row := r.Db.QueryRow(`
		SELECT * FROM 
			profiles 
		WHERE 
			id = $1 and deleted_at is null`, id)

	err = row.Scan(
		&profile.ID,
		&profile.FullName,
		&profile.CountryCode,
		&profile.PhoneNumber,
		&profile.Password,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.DeletedAt,
	)
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (r *Repository) UpdateProfileByID(profile Profile) (err error) {
	query := "UPDATE profiles SET"
	setValues := make([]string, 0)

	if profile.FullName != "" {
		setValues = append(setValues, fmt.Sprintf("full_name = '%s'", profile.FullName))
	}
	if profile.PhoneNumber != "" {
		setValues = append(setValues, fmt.Sprintf("phone_number = '%s'", profile.PhoneNumber))
	}

	// If there are fields to update, execute the query
	if len(setValues) > 0 {
		query += " " + strings.Join(setValues, ", ") + " WHERE id = $1"
		_, err := r.Db.Exec(query, profile.ID)
		return err
	}

	// No fields to update
	return nil
}

func (r *Repository) UpsertProfileMetaData(input ProfileMetaData) (createdID int, err error) {

	// This method will increment login_attempt by 1
	// If the corresponding row has not been created yet, it will insert a new row with initial login attempt 1
	// If the corresponding row has been created, it will update the login attempt

	stmt, err := r.Db.Prepare(
		`INSERT INTO profile_metadata
			(
				profile_id, 
				login_attempt
			) VALUES ($1, 1) 
		ON CONFLICT (profile_id)
			DO 
				UPDATE 
					SET login_attempt = profile_metadata.login_attempt + 1, updated_at = $2
		RETURNING id`)

	if err != nil {
		return createdID, err
	}

	err = stmt.QueryRow(
		input.ProfileID,
		time.Now(),
	).Scan(&createdID)

	if err != nil {
		return createdID, err
	}

	return createdID, nil
}
