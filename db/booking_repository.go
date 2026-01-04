package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"tickets/entities"
	"tickets/message/event"
	"tickets/message/outbox"
)

type BookingsRepository struct {
	db *sqlx.DB
}

func NewBookingsRepository(db *sqlx.DB) BookingsRepository {
	if db == nil {
		panic("nil db")
	}

	return BookingsRepository{db: db}
}

func (b BookingsRepository) AddBooking(ctx context.Context, booking entities.Booking) (err error) {
	return updateInTx(
		ctx,
		b.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			_, err = tx.NamedExecContext(ctx, `
				INSERT INTO
					bookings (booking_id, show_id, number_of_tickets, customer_email)
				VALUES (:booking_id, :show_id, :number_of_tickets, :customer_email)
		`, booking)
			if err != nil {
				return fmt.Errorf("could not add booking: %w", err)
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			err = event.NewBus(outboxPublisher).Publish(ctx, entities.BookingMade{
				Header:          entities.NewMessageHeader(),
				BookingID:       booking.BookingID,
				NumberOfTickets: booking.NumberOfTickets,
				CustomerEmail:   booking.CustomerEmail,
				ShowID:          booking.ShowID,
			})
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}

			return nil
		},
	)
}
