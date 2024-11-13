package account

type AccountType string

const (
	TypeSending   AccountType = "sending"
	TypeReceiving AccountType = "receiving"
)

type Account struct {
	// number is unique identifier of acc
	Number  uint        `json:"number"`
	Name    string      `json:"name"`
	Iban    string      `json:"iban"`
	Address string      `json:"address"`
	Amount  uint        `json:"amount"`
	Type    AccountType `json:"type"`
}
