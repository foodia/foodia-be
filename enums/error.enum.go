package enums

import "errors"

var (
	ErrAccessForbidden          = errors.New("you are not authorized to acces this url")
	ErrUnauthorized             = errors.New("Unauthorized")
	ErrNotFound                 = errors.New("Requested data is not found")
	ErrInternalServor           = errors.New("Internal Server Error")
	ErrBadParamInput            = errors.New("Requested parameters are not valid")
	ErrIncorrectCredential      = errors.New("Login failed. Email or password is incorrect.")
	ErrInvalidToken             = errors.New("token is invalid")
	ErrInvalidRefreshToken      = errors.New("refresh token is invalid")
	ErrExpiredToken             = errors.New("token has expired")
	ErrEmailOrPasswordMissMatch = errors.New("email/password miss match")
)
