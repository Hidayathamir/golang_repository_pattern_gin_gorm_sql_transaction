package model

type User struct {
	ID       int    `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	FullName string `gorm:"column:full_name"`
	Balance  int    `gorm:"column:balance"`
}

type Transaction struct {
	ID          int `gorm:"column:id"`
	SenderID    int `gorm:"column:sender_id"`
	RecipientID int `gorm:"column:recipient_id"`
	Amount      int `gorm:"column:amount"`
}
