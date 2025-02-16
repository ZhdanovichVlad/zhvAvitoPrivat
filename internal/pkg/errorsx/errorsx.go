package errorsx

import "fmt"

var (
	ErrInvalidPassword      = fmt.Errorf("invalid password")
	ErrDB                   = fmt.Errorf("database error")
	ErrService              = fmt.Errorf("service error")
	ErrInvalidInput         = fmt.Errorf("invalid input")
	ErrUnknown              = fmt.Errorf("unknow error")
	ErrTokenExpired         = fmt.Errorf("token expired")
	ErrWrongUUID            = fmt.Errorf("incorrect uuid")
	ErrAuthHeaderIsEmpty    = fmt.Errorf("requiredAuthorization header is required")
	ErrUnexpSignedMetod     = fmt.Errorf("unexpected signing method")
	ErrInvUserUUIDInToken   = fmt.Errorf("invalid userUUID in token")
	ErrInvalidToken         = fmt.Errorf("invalid token")
	ErrUnknownUser          = fmt.Errorf("unknown user")
	ErrItemNotFound         = fmt.Errorf("item not found")
	ErrNotEnoughMoney       = fmt.Errorf("not enough money")
	ErrReceiverNotFound     = fmt.Errorf("receiver Not Found")
	ErrForbiddenTransaction = fmt.Errorf("you can't transfer money to yourself")
)
