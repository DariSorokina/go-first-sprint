package customerrors

import "errors"

// ErrShortURLAlreadyExist indicates that the corresponding short URL already exists.
var ErrShortURLAlreadyExist = errors.New("corresponding short URL already exists")
