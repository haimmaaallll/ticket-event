package exception

import "errors"

var (
	ErrEmailTaken        = errors.New("Email sudah terdaftar")
	ErrInvalidCredential = errors.New("email atau password salah")
)
