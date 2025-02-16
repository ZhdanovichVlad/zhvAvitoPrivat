package entity

type UserInfo struct {
	Coins       int         `json:"coins"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coin_history"`
}

type Item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity,omitempty"`
}

type CoinHistory struct {
	Received []UserTransfer `json:"received"`
	Sent     []UserTransfer `json:"sent"`
}

type UserTransfer struct {
	User   string `json:"user"`
	Amount int    `json:"amount"`
}

type Transaction struct {
	Sender       string
	SenderUUID   string
	Receiver     string
	ReceiverUUID string
	Amount       int
}
