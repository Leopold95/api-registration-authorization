package domain

import "errors"

var ErrorUserExists = errors.New("user already exists")
var ErrorRegistrationConflict = errors.New("registration id belongs to different user data")
