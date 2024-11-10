package account

type AccountType string

const (
	TypeSending   AccountType = "sending"
	TypeReceiving AccountType = "receiving"
)

type Account struct {
	// number is unique identifier of acc
	Number  uint        `json:"number,omitempty"`
	Name    string      `json:"name,omitempty"`
	Iban    string      `json:"iban,omitempty"`
	Address string      `json:"address,omitempty"`
	Amount  uint        `json:"amount,omitempty"`
	Type    AccountType `json:"type,omitempty"`
}
