package handlers

import (
	"net/http"
	"time"
	"txsystem/internal/ledger/models"
	"txsystem/internal/ledger/service"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type LedgerHandler struct {
	service *service.LedgerService
}

func NewLedgerHandler(db *mongo.Database) *LedgerHandler {
	return &LedgerHandler{
		service: service.NewLedgerService(db),
	}
}

// ListLedgersByAccount retrieves all ledger entries for a specific account
func (h *LedgerHandler) ListLedgersByAccount(c echo.Context) error {
	accountID := c.Param("accountId")
	if accountID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "account ID is required",
		})
	}

	ctx := c.Request().Context()
	ledgers, err := h.service.ListLedgers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch ledger entries",
		})
	}

	// Filter by account ID
	var accountLedgers []*models.Ledger
	for _, ledger := range ledgers {
		if ledger.AccountID == accountID {
			accountLedgers = append(accountLedgers, ledger)
		}
	}

	return c.JSON(http.StatusOK, accountLedgers)
}

func (h *LedgerHandler) ListAllLedgersByDate(c echo.Context) error {
	dateStr := c.QueryParam("date")

	var filterDate time.Time
	var err error

	if dateStr != "" {
		filterDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid date format, please use YYYY-MM-DD",
			})
		}
	}

	ctx := c.Request().Context()
	ledgers, err := h.service.ListLedgers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch ledger entries",
		})
	}

	if dateStr == "" {
		return c.JSON(http.StatusOK, ledgers)
	}

	var filteredLedgers []*models.Ledger
	for _, ledger := range ledgers {
		ledgerDate := ledger.CreatedAt.Truncate(24 * time.Hour)
		filterDateOnly := filterDate.Truncate(24 * time.Hour)

		if ledgerDate.Equal(filterDateOnly) {
			filteredLedgers = append(filteredLedgers, ledger)
		}
	}

	return c.JSON(http.StatusOK, filteredLedgers)
}

func InitRoutes(e *echo.Echo, db *mongo.Database) {
	h := NewLedgerHandler(db)

	e.Logger.Info("Initializing ledger routes")
	g := e.Group("/api/v1/ledger")
	g.GET("/account/:accountId", h.ListLedgersByAccount)
	g.GET("/", h.ListAllLedgersByDate) // New endpoint to get all ledgers with optional date filtering

}
