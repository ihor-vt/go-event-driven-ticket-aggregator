package entities

type Booking struct {
	BookingID       string `json:"booking_id" db:"booking_id"`
	ShowID          string `json:"show_id" db:"show_id"`
	NumberOfTickets int    `json:"number_of_tickets" db:"number_of_tickets"`
	CustomerEmail   string `json:"customer_email" db:"customer_email"`
}

type BookingMade struct {
	Header MessageHeader `json:"header"`

	NumberOfTickets int    `json:"number_of_tickets"`
	BookingID       string `json:"booking_id"`
	CustomerEmail   string `json:"customer_email"`
	ShowID          string `json:"show_id"`
}
