package types

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusCompleted TransactionStatus = "completed"
	StatusFailed    TransactionStatus = "failed"
	// add more as needed
)

type TransactionRequest struct {
	Amount             float64 `json:"amount"`
	Description        string  `json:"description"`
	SourceAccount      string  `json:"source_account"`
	DestinationAccount string  `json:"destination_account"`
	TransactionType    string  `json:"transaction_type"`
}

type TransactionResponse struct {
	ID                 uint64  `json:"id"`
	Amount             float64 `json:"amount"`
	Description        string  `json:"description"`
	SourceAccount      string  `json:"source_account"`
	DestinationAccount string  `json:"destination_account"`
	TransactionType    string  `json:"transaction_type"`
	Status             string  `json:"status"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
	TransactionID      string  `json:"transaction_id"`
}
