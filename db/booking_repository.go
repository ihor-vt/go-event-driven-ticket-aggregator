package db

import (
	"context"
	"database/sql"
	"fmt"
	"tickets/entities"
	"tickets/message/event"
	"tickets/message/outbox"

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

func (s BookingRepository) AddBooking(ctx context.Context, booking entities.Booking) (err error) {
	return updateInTx(
		ctx,
		s.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			_, err = tx.NamedExecContext(ctx, `
				INSERT INTO
					bookings (booking_id, show_id, number_of_tickets, customer_email)
				VALUES (:booking_id, :show_id, :number_of_tickets, :customer_email)
			`, booking)

			outboxPublish, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}
			bus := event.NewBus(outboxPublish)

			err = bus.Publish(ctx, entities.BookingMade{
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
