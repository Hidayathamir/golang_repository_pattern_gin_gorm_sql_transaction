package repository

import (
	"context"

	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/database/model"
	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/util/transaction"
	"gorm.io/gorm"
)

type IPaymentRepo interface {
	GetUserByID(ctx context.Context, ID int) (model.User, error)
	CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
	UpdateUserBalance(ctx context.Context, userID int, balance int) error
}

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) IPaymentRepo {
	return &PaymentRepo{db: db}
}

func (p *PaymentRepo) getTxOrDB(ctx context.Context) *gorm.DB {
	isHasTransaction := ctx.Value(transaction.CtxKey) != nil
	if isHasTransaction {
		if tx, ok := ctx.Value(transaction.CtxKey).(*gorm.DB); ok {
			return tx
		}
	}
	return p.db
}

func (p *PaymentRepo) GetUserByID(ctx context.Context, ID int) (model.User, error) {
	user := model.User{}
	qx := p.getTxOrDB(ctx).First(&user, ID)
	if qx.Error != nil {
		return model.User{}, qx.Error
	}
	return user, nil
}

func (p *PaymentRepo) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	qx := p.getTxOrDB(ctx).Create(&transaction)
	if qx.Error != nil {
		return model.Transaction{}, qx.Error
	}
	return transaction, nil
}

func (p *PaymentRepo) UpdateUserBalance(ctx context.Context, userID int, balance int) error {
	return p.getTxOrDB(ctx).
		Table("users").
		Where("id = ?", userID).
		Update("balance", balance).Error
}
