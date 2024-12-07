package routes

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/denys89/ewallet-api/config"
	"github.com/denys89/ewallet-api/middleware"
	"github.com/denys89/ewallet-api/models"
	"github.com/denys89/ewallet-api/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	RecipientID string  `json:"target_user,omitempty"`
	Description string  `json:"remarks,omitempty"`
}

type PaymentRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"remarks" binding:"required"`
}

func TopUp(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactionRepo := repositories.NewTransactionRepository(config.DB)
	transactionID, balanceBefore, balanceAfter, err := transactionRepo.TopUp(userID, req.Amount)
	if err != nil {
		log.Printf("Top-up error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process top-up"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"top_up_id":      transactionID,
			"amount_top_up":  req.Amount,
			"balance_before": balanceBefore,
			"balance_after":  balanceAfter,
			"created_date":   time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

func Transfer(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RecipientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target user is required"})
		return
	}

	transactionRepo := repositories.NewTransactionRepository(config.DB)
	transaction, balanceBefore, balanceAfter, err := transactionRepo.Transfer(userID, req.Amount, req.RecipientID, req.Description)
	if err != nil {
		log.Printf("Transfer error: %v", err)
		if err == models.ErrInvalidTransaction {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Balance is not enough"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"transfer_id":    transaction.ID,
			"amount":         transaction.Amount,
			"balance_before": balanceBefore,
			"balance_after":  balanceAfter,
			"remarks":        transaction.Description,
			"created_date":   transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func Payment(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactionRepo := repositories.NewTransactionRepository(config.DB)
	transaction, balanceBefore, balanceAfter, err := transactionRepo.Payment(userID, req.Amount, req.Description)
	if err != nil {
		log.Printf("Payment error: %v", err)
		if err == models.ErrInvalidTransaction {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Balance is not enough"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"payment_id":     transaction.ID,
			"amount":         transaction.Amount,
			"balance_before": balanceBefore,
			"balance_after":  balanceAfter,
			"remark":         transaction.Description,
			"created_at":     transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func GetTransactionHistory(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	transactionRepo := repositories.NewTransactionRepository(config.DB)
	transactions, err := transactionRepo.GetUserTransactions(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	var transactionResponses []gin.H
	for _, t := range transactions {
		transactionResponses = append(transactionResponses, gin.H{
			"id":             t.ID,
			"amount":         t.Amount,
			"type":           t.Type,
			"remarks":        t.Description,
			"balance_before": t.BalanceBefore,
			"balance_after":  t.BalanceAfter,
			"status":         t.Status,
			"created_date":   t.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result": transactionResponses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
		},
	})
}
