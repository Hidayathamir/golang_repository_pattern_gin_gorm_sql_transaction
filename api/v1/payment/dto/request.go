package dto

type ReqTransfer struct {
	SenderID    int `json:"sender_id"`
	RecipientID int `json:"recipient_id"`
	Amount      int `json:"amount"`
}
