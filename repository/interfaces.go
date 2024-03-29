// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

type RepositoryInterface interface {
	GetPhoneNumberExistence(phoneNumber string) (isExist bool, err error)
	GetPhoneNumberExistenceWithExcludedID(phoneNumber string, excludedID int) (isExist bool, err error)
	GetProfileByPhoneNumber(phoneNumber string) (profile Profile, err error)
	GetProfileByID(id int) (profile Profile, err error)
	CreateProfile(input Profile) (createdID int, err error)
	UpdateProfileByID(profile Profile) (err error)
	UpsertProfileMetaData(input ProfileMetaData) (createdID int, err error)
}
