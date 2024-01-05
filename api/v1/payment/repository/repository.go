package repository

import (
	"context"

	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/database/model"
	"github.com/Hidayathamir/txmanager"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IPaymentRepo interface {
	GetUserByID(ctx context.Context, ID int) (model.User, error)
	CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
	UpdateUserBalance(ctx context.Context, userID int, balance int) error
}

type PaymentRepo struct {
	db        *gorm.DB
	txManager txmanager.ITransactionManager
}

func NewPaymentRepo(db *gorm.DB, txManager txmanager.ITransactionManager) IPaymentRepo {
	return &PaymentRepo{db: db, txManager: txManager}
}

func (p *PaymentRepo) GetUserByID(ctx context.Context, ID int) (model.User, error) {
	db := p.db

	if tx, ok := p.txManager.GetTx(ctx); ok {
		db = tx.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	user := model.User{}
	err := db.First(&user, ID).Error
	return user, err
}

func (p *PaymentRepo) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	db := p.db

	if tx, ok := p.txManager.GetTx(ctx); ok {
		db = tx
	}

	err := db.Create(&transaction).Error
	return transaction, err
}

func (p *PaymentRepo) UpdateUserBalance(ctx context.Context, userID int, balance int) error {
	db := p.db

	if tx, ok := p.txManager.GetTx(ctx); ok {
		db = tx
	}

	err := db.Table("users").Where("id = ?", userID).Update("balance", balance).Error
	return err
}
