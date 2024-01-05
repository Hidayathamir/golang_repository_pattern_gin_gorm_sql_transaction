# golang_repository_pattern_gin_gorm_sql_transaction
Golang repository pattern gin gorm sql transaction

## How To Transaction
Check file [service.go](./api/v1/payment/service/service.go). `Transfer` method do `txManager.SQLTransaction`, it can be nested like when it's called `p.updateBalanceSenderAndRecipient` which have it's own `txManager.SQLTransaction`.

```go
package service

import (
	"context"
	"errors"

	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/api/v1/payment/dto"
	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/api/v1/payment/repository"
	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/database/model"
	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/util/transaction"
)

type IPaymentService interface {
	Transfer(ctx context.Context, req dto.ReqTransfer) (model.Transaction, error)
}

type PaymentService struct {
	repo      repository.IPaymentRepo
	txManager transaction.ITransactionManager
}

func NewPaymentService(repo repository.IPaymentRepo, txManager transaction.ITransactionManager) IPaymentService {
	return &PaymentService{repo: repo, txManager: txManager}
}

func (p *PaymentService) Transfer(ctx context.Context, req dto.ReqTransfer) (model.Transaction, error) {
	if err := validateReqTransfer(req); err != nil {
		return model.Transaction{}, err
	}

	var transaction model.Transaction
	err := p.txManager.SQLTransaction(ctx, func(ctx context.Context) error {
		sender, err := p.repo.GetUserByID(ctx, req.SenderID)
		if err != nil {
			return err
		}

		if req.Amount > sender.Balance {
			return errors.New("balance is not enough")
		}

		recipient, err := p.repo.GetUserByID(ctx, req.RecipientID)
		if err != nil {
			return err
		}

		transaction, err = p.repo.CreateTransaction(ctx, model.Transaction{
			SenderID:    req.SenderID,
			RecipientID: req.RecipientID,
			Amount:      req.Amount,
		})
		if err != nil {
			return err
		}

		err = p.updateBalanceSenderAndRecipient(ctx, req.Amount, sender, recipient)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}

func validateReqTransfer(req dto.ReqTransfer) error {
	if req.SenderID == 0 {
		return errors.New("sender_id can not be empty")
	}

	if req.RecipientID == 0 {
		return errors.New("recipient_id can not be empty")
	}

	if req.SenderID == req.RecipientID {
		return errors.New("can not transfer to yourself")
	}

	if req.Amount <= 10000 {
		return errors.New("amount can not be less than 10000")
	}

	return nil
}

func (p *PaymentService) updateBalanceSenderAndRecipient(ctx context.Context, transferAmount int, sender model.User, recipient model.User) error {
	// txManager.SQLTransaction can be nested
	return p.txManager.SQLTransaction(ctx, func(ctx context.Context) error {
		err := p.repo.UpdateUserBalance(ctx, sender.ID, sender.Balance-transferAmount)
		if err != nil {
			return err
		}

		err = p.repo.UpdateUserBalance(ctx, recipient.ID, recipient.Balance+transferAmount)
		if err != nil {
			return err
		}

		return nil
	})
}
```

## Quick Start

1. Run MySQL database using Docker Compose:

```shell
sudo docker compose up
```

MySQL credentials is:

```
DB_USER="user"
DB_PASSWORD="password"
DB_HOST="localhost"
DB_PORT="9306"
DB_NAME="bank_central_jakarta"
```

2. Run go app:

```shell
go run .
```

3. Hit Api Transfer:

```shell
curl --location 'localhost:8080/api/v1/payment/transfer' \
--header 'Content-Type: application/json' \
--data '{
    "sender_id":3,
    "recipient_id":1,
    "amount":50000
}'
```
