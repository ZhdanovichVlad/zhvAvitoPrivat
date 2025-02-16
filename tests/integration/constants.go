package integration

const (
	findUserPath                = "SELECT uuid, username, balance FROM users WHERE username = $1"
	findMerchInInventory        = "SELECT merchandise_uuid, quantity FROM owned_inventory WHERE  user_uuid = $1"
	findMerchPrice              = "SELECT uuid, price FROM merchandise WHERE name = $1"
	findTransactionsBySender    = "SELECT sender_uuid, recipient_uuid, quantity FROM transactions WHERE sender_uuid = $1"
	findTransactionsByRecipient = "SELECT sender_uuid, recipient_uuid, quantity FROM transactions WHERE recipient_uuid = $1"
)
