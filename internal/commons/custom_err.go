package commons

import "errors"

var (
	ErrUserRequired        = errors.New("User is required")
	ErrPasswordRequired    = errors.New("Password is required")
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound            = errors.New("Your requested Item is not found")
	ErrConflict            = errors.New("Your Item already exist")
	ErrBadParamInput       = errors.New("Given Param is not valid")
	ErrCannotBeDeleted     = errors.New("This Item can't be deleted")
	ErrInvalidCredentials  = errors.New("Invalid Credentials")
)
