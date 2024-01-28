package repository

import (
	"database/sql"
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
