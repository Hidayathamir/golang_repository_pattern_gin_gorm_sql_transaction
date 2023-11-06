package controller

import (
	"net/http"

	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/api/v1/payment/dto"
	"github.com/Hidayathamir/golang_repository_pattern_gin_gorm_sql_transaction/api/v1/payment/service"
	"github.com/gin-gonic/gin"
)

type IPaymentController interface {
	Transfer(ctx *gin.Context)
}

type PaymentController struct {
	service service.IPaymentService
}

func NewPaymentController(service service.IPaymentService) IPaymentController {
	return &PaymentController{service: service}
}

func (p *PaymentController) Transfer(ctx *gin.Context) {
	req := dto.ReqTransfer{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	transfer, err := p.service.Transfer(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  transfer.ID,
		"error": nil,
	})
}
