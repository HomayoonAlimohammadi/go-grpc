package server

import "errors"

var ErrorLaptopAlreadyExists = errors.New("laptop already exists")

var ErrorCanNotCopyLaptop = errors.New("can not copy laptop")
