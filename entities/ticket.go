package entities

type Ticket struct {
	TicketID      string `json:"ticket_id" db:"ticket_id"`
	Price         Money  `json:"price" db:"price"`
	CustomerEmail string `json:"customer_email" db:"customer_email"`
}

type TicketPrinted struct {
	Header MessageHeader `json:"header"`

	TicketID string `json:"ticket_id"`
	FileName string `json:"file_name"`
}
