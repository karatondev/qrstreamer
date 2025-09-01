package model

type WhatsappAccount struct {
	AccountID     string  `json:"account_id"` // UUID
	UserID        string  `json:"user_id"`    // UUID
	AccountName   string  `json:"account_name"`
	AccountAlias  *string `json:"account_alias,omitempty"`
	PhoneNumber   *string `json:"phone_number,omitempty"`
	SenderJID     *string `json:"sender_jid,omitempty"`
	ConnectStatus string  `json:"connect_status"`
	IsActive      bool    `json:"is_active"`
}
