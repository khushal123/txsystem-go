package handler

import (
	"strconv"
	"txsystem/internal/account/service"
	"txsystem/internal/common/types"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	service *service.AccountService
}

func NewHandler(s *service.AccountService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) GetAccount(c echo.Context) error {
	id := c.Param("id")
	accountID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid account ID"})
	}
	account, err := h.service.GetAccount(c.Request().Context(), accountID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to get account"})
	}
	return c.JSON(200, account)
}

func InitRoutes(e *echo.Echo, kc types.ProducerConnection, db *gorm.DB) {
	transactionService := service.NewAccountService(db)
	h := NewHandler(transactionService)
	e.Logger.Info("Initializing transaction routes")
	g := e.Group("/api/v1/transactions")
	g.GET("/:id", h.GetAccount)
}
