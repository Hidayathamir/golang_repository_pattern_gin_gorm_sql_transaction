package repository

import (
	"context"

	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/database/model"
	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/util/transaction"
)

type IPaymentRepo interface {
	GetUserByID(ctx context.Context, ID int) (model.User, error)
	CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
	UpdateUserBalance(ctx context.Context, userID int, balance int) error
}

type PaymentRepo struct {
	txManager transaction.ITransactionManager
}

func NewPaymentRepo(txManager transaction.ITransactionManager) IPaymentRepo {
	return &PaymentRepo{txManager: txManager}
}

func (p *PaymentRepo) GetUserByID(ctx context.Context, ID int) (model.User, error) {
	user := model.User{}
	qx := p.txManager.GetTxOrDB(ctx).First(&user, ID)
	if qx.Error != nil {
		return model.User{}, qx.Error
	}
	return user, nil
}

func (p *PaymentRepo) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	qx := p.txManager.GetTxOrDB(ctx).Create(&transaction)
	if qx.Error != nil {
		return model.Transaction{}, qx.Error
	}
	return transaction, nil
}

func (p *PaymentRepo) UpdateUserBalance(ctx context.Context, userID int, balance int) error {
	return p.txManager.GetTxOrDB(ctx).
		Table("users").
		Where("id = ?", userID).
		Update("balance", balance).Error
}
