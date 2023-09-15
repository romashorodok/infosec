package identity

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type TokenResult struct {
	AccessToken  string
	RefreshToken http.Cookie
}

type userInfo struct {
	userID   uuid.UUID
	username string
	claims   []string
}

var (
	NotFoundUserError       = errors.New("Not found user")
	InvalidCredentialsError = errors.New("Invalid credentials")
	TransactionFailedError  = errors.New("Unable start transaction")
	UnableHashPassword      = errors.New("Unable generate hash of password")
	UnableStoreUser         = errors.New("Unable store user record")
)
