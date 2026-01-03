package db

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type BookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) BookingRepository {
	if db == nil {
		panic("db is nil")
	}

	return BookingRepository{db: db}
}

func (s BookingRepository) AddBooking(ctx context.Context, booking entities.Booking) error {
	_, err := s.db.NamedExecContext(ctx, `
		INSERT INTO
			bookings (booking_id, show_id, number_of_tickets, customer_email)
		VALUES (:booking_id, :show_id, :number_of_tickets, :customer_email)
	`, booking)
	if err != nil {
		return fmt.Errorf("could not add booking: %w", err)
	}

	return nil
}
