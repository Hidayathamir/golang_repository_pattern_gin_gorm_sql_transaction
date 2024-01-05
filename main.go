package main

import (
	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/api/v1/payment/controller"
	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/api/v1/payment/repository"
	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/api/v1/payment/service"
	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/util/transaction"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := getDBConnection()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	{
		paymentController := getPaymentController(db)
		r.POST("/api/v1/payment/transfer", paymentController.Transfer)
	}

	if err := r.Run(); err != nil {
		panic(err)
	}
}

func getDBConnection() (*gorm.DB, error) {
	gormConfig := gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	dsn := "user:password@tcp(localhost:9306)/bank_central_jakarta?parseTime=true"
	return gorm.Open(mysql.Open(dsn), &gormConfig)
}

func getPaymentController(db *gorm.DB) controller.IPaymentController {
	txManager := transaction.NewTransactionManager(db)
	repo := repository.NewPaymentRepo(db, txManager)
	service := service.NewPaymentService(repo, txManager)
	return controller.NewPaymentController(service)
}
