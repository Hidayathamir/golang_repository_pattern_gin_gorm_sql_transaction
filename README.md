# golang_repository_pattern_gin_gorm_sql_transaction
Golang repository pattern gin gorm sql transaction

## How To Transaction
Check file [service.go](./api/v1/payment/service/service.go). `Transfer` method do `txManager.SQLTransaction`, it can be nested like when it's called `p.updateBalanceSenderAndRecipient` which have it's own `txManager.SQLTransaction`.

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
