package customErrors

import "errors"

var ShortURLAlreadyExistError = errors.New("Corresponding short URL already exists")
