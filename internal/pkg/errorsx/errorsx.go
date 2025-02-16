package errorsx

import "fmt"

var (
	InvalidPassword      = fmt.Errorf("invalid password")
	DBError              = fmt.Errorf("database error")
	ServiceError         = fmt.Errorf("service error")
	InvalidInput         = fmt.Errorf("invalid input")
	UnknownError         = fmt.Errorf("unknow error")
	TokenExpired         = fmt.Errorf("token expired")
	Unauthorized         = fmt.Errorf("unauthorized")
	AuthHeaderIsEmpty    = fmt.Errorf("requiredAuthorization header is required")
	UnexpSignedMetod     = fmt.Errorf("unexpected signing method")
	InvUserUuidInToken   = fmt.Errorf("invalid userUuid in token")
	InvalidToken         = fmt.Errorf("invalid token")
	UnknownUser          = fmt.Errorf("unknown user")
	ItemNotFound         = fmt.Errorf("item not found")
	NotEnoughMoney       = fmt.Errorf("not enough money")
	ReceiverNotFound     = fmt.Errorf("receiver Not Found")
	ForbiddenTransaction = fmt.Errorf("you can't transfer money to yourself")
)
