package models

import "errors"

var (
	ErrInvalidZipCode  = errors.New("invalid zip code")
	ErrZipCodeNotFound = errors.New("zip code not found")
	ErrInternalServer  = errors.New("internal server error")
	ErrWeatherNotFound = errors.New("weather not found")
)
