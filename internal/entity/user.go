package entity

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Uuid         string
	Username     string
	PasswordHash string
	Balance      int
}

type SendingCoins struct {
	User   string `json:"toUser"`
	Amount int    `json:"amount"`
}
