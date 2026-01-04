package entities

type VoidReceipt struct {
	TicketID       string
	Reason         string
	IdempotencyKey string
}

type IssueReceiptRequest struct {
	TicketID string `json:"ticket_id"`
	Price    Money  `json:"price"`

	IdempotencyKey string `json:"idempotency_key"`
}
