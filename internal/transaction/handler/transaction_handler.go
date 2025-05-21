package handler

import (
	"net/http"
	"strconv"

	"txsystem/internal/transaction/repository"
	"txsystem/internal/transaction/service"
	"txsystem/pkg/common/types"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	service *service.TransactionService
}

func NewHandler(s *service.TransactionService) *Handler {
	return &Handler{
		service: s,
	}
}

// @Summary Create a new transaction
// @Description CreateTransaction handles creating a transaction.
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body types.TransactionRequest true "Transaction request"
// @Success 201 {object} map[string]string "message:transaction created"
// @Failure 400 {object} map[string]string "error:invalid request"
// @Failure 500 {object} map[string]string "error:failed to create transaction"
// @Router /api/v1/transactions [post]
func (h *Handler) CreateTransaction(c echo.Context) error {
	var req types.TransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := h.service.CreateTransaction(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create transaction"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "transaction created"})
}

// @Summary Get transactions
// @Description GetTransactions handles fetching list of last 10 transactions
// @Tags transactions
// @Produce json
// @Success 200 {array} types.TransactionResponse "List of transactions"
// @Failure 400 {object} map[string]string "error:bad request"
// @Failure 404 {object} map[string]string "error:transactions not found"
// @Failure 500 {object} map[string]string "error:failed to fetch transactions"
// @Router /api/v1/transactions [get]
func (h *Handler) GetTransactions(c echo.Context) error {
	transactions, err := h.service.GetTransactions(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch transactions"})
	}

	return c.JSON(http.StatusOK, transactions)
}

// @Summary Get a transaction by ID
// @Description GetTransaction handles fetching a transaction by ID.
// @Tags transactions
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} types.TransactionResponse "Transaction details"
// @Failure 400 {object} map[string]string "error:invalid transaction ID"
// @Failure 404 {object} map[string]string "error:transaction not found"
// @Failure 500 {object} map[string]string "error:failed to fetch transaction"
// @Router /api/v1/transactions/{id} [get]
func (h *Handler) GetTransaction(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid transaction ID"})
	}

	tx, err := h.service.GetTransaction(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch transaction"})
	}
	if tx == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "transaction not found"})
	}

	return c.JSON(http.StatusOK, tx)
}

func InitRoutes(e *echo.Echo, kc types.ProducerConnection, db *gorm.DB) {
	transactionService := service.NewTransactionService(kc, repository.NewTransactionRepository(db))
	h := NewHandler(transactionService)
	e.Logger.Info("Initializing transaction routes")
	g := e.Group("/api/v1/transactions")
	g.POST("", h.CreateTransaction)
	g.GET("", h.GetTransactions)
	g.GET("/:id", h.GetTransaction)
}
