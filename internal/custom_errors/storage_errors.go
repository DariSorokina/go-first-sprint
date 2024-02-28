package customerrors

import "errors"

var ErrShortURLAlreadyExist = errors.New("corresponding short URL already exists")
