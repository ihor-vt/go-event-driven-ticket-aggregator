package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"tickets/entities"
)

type OpsBookingReadModel struct {
	db *sqlx.DB
}

func NewOpsBookingReadModel(db *sqlx.DB) OpsBookingReadModel {
	if db == nil {
		panic("nil db")
	}
	return OpsBookingReadModel{db: db}
}

func (r OpsBookingReadModel) OnBookingMade(ctx context.Context, event *entities.BookingMade) error {
	readModel := entities.OpsBooking{
		BookingID:  event.BookingID,
		BookedAt:   event.Header.PublishedAt,
		Tickets:    make(map[string]entities.OpsTicket),
		LastUpdate: time.Now(),
	}

	return r.createReadModel(ctx, readModel)
}

func (r OpsBookingReadModel) OnTicketBookingConfirmed(
	ctx context.Context,
	event *entities.TicketBookingConfirmed,
) error {
	return r.updateReadModelByBookingID(
		ctx,
		event.BookingID,
		func(rm entities.OpsBooking) (entities.OpsBooking, error) {
			ticket := rm.Tickets[event.TicketID]

			ticket.PriceAmount = event.Price.Amount
			ticket.PriceCurrency = event.Price.Currency
			ticket.CustomerEmail = event.CustomerEmail
			ticket.Status = "confirmed"

			rm.Tickets[event.TicketID] = ticket

			return rm, nil
		},
	)
}

func (r OpsBookingReadModel) OnTicketReceiptIssued(
	ctx context.Context,
	event *entities.TicketReceiptIssued,
) error {
	return r.updateReadModelByTicketID(
		ctx,
		event.TicketID,
		func(ticket entities.OpsTicket) (entities.OpsTicket, error) {
			ticket.ReceiptIssuedAt = event.IssuedAt
			ticket.ReceiptNumber = event.ReceiptNumber

			return ticket, nil
		},
	)
}

func (r OpsBookingReadModel) OnTicketPrinted(
	ctx context.Context,
	event *entities.TicketPrinted,
) error {
	return r.updateReadModelByTicketID(
		ctx,
		event.TicketID,
		func(ticket entities.OpsTicket) (entities.OpsTicket, error) {
			ticket.PrintedAt = event.Header.PublishedAt
			ticket.PrintedFileName = event.FileName

			return ticket, nil
		},
	)
}

func (r OpsBookingReadModel) OnTicketRefunded(
	ctx context.Context,
	event *entities.TicketRefunded,
) error {
	return r.updateReadModelByTicketID(
		ctx,
		event.TicketID,
		func(ticket entities.OpsTicket) (entities.OpsTicket, error) {
			ticket.Status = "refunded"

			return ticket, nil
		},
	)
}

func (r OpsBookingReadModel) AllBookings(ctx context.Context, receiptIssueDate *string) ([]entities.OpsBooking, error) {
	query := "SELECT payload FROM read_model_ops_bookings"
	var queryArgs []any

	if receiptIssueDate != nil {
		query += `
			WHERE booking_id IN (
				SELECT booking_id FROM (
					SELECT booking_id, 
						DATE(jsonb_path_query(payload, '$.tickets.*.receipt_issued_at')::text) as receipt_issued_at 
					FROM 
						read_model_ops_bookings
				) bookings_within_date 
				WHERE receipt_issued_at = $1
			)`
		queryArgs = append(queryArgs, *receiptIssueDate)
	}

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.OpsBooking
	for rows.Next() {
		var payload []byte
		if err := rows.Scan(&payload); err != nil {
			return nil, err
		}

		reservation, err := r.unmarshalReadModelFromDB(payload)
		if err != nil {
			return nil, err
		}

		result = append(result, reservation)
	}

	return result, nil
}

func (r OpsBookingReadModel) BookingReadModel(
	ctx context.Context,
	bookingID string,
) (entities.OpsBooking, error) {
	return r.findReadModelByBookingID(ctx, bookingID, r.db)
}

func (r OpsBookingReadModel) createReadModel(
	ctx context.Context,
	booking entities.OpsBooking,
) (err error) {
	payload, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO 
		    read_model_ops_bookings (payload, booking_id)
		VALUES
			($1, $2)
		ON CONFLICT (booking_id) DO NOTHING;
`, payload, booking.BookingID)
	if err != nil {
		return fmt.Errorf("could not create read model: %w", err)
	}

	return nil
}

func (r OpsBookingReadModel) updateReadModelByBookingID(
	ctx context.Context,
	bookingID string,
	updateFunc func(rm entities.OpsBooking) (entities.OpsBooking, error),
) (err error) {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			rm, err := r.findReadModelByBookingID(ctx, bookingID, tx)
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("read model for booking %s not exist yet", bookingID)
			} else if err != nil {
				return fmt.Errorf("could not find read model: %w", err)
			}

			updatedRm, err := updateFunc(rm)
			if err != nil {
				return err
			}

			return r.updateReadModel(ctx, tx, updatedRm)
		},
	)
}

func (r OpsBookingReadModel) updateReadModelByTicketID(
	ctx context.Context,
	ticketID string,
	updateFunc func(ticket entities.OpsTicket) (entities.OpsTicket, error),
) (err error) {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			rm, err := r.findReadModelByTicketID(ctx, ticketID, tx)
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("read model for ticket %s not exist yet", ticketID)
			} else if err != nil {
				return fmt.Errorf("could not find read model: %w", err)
			}

			ticket, _ := rm.Tickets[ticketID]

			updatedRm, err := updateFunc(ticket)
			if err != nil {
				return err
			}

			rm.Tickets[ticketID] = updatedRm

			return r.updateReadModel(ctx, tx, rm)
		},
	)
}

func (r OpsBookingReadModel) updateReadModel(
	ctx context.Context,
	tx *sqlx.Tx,
	rm entities.OpsBooking,
) error {
	rm.LastUpdate = time.Now()

	payload, err := json.Marshal(rm)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO 
			read_model_ops_bookings (payload, booking_id)
		VALUES
			($1, $2)
		ON CONFLICT (booking_id) DO UPDATE SET payload = excluded.payload;
		`, payload, rm.BookingID)
	if err != nil {
		return fmt.Errorf("could not update read model: %w", err)
	}

	return nil
}

func (r OpsBookingReadModel) findReadModelByTicketID(
	ctx context.Context,
	ticketID string,
	db dbExecutor,
) (entities.OpsBooking, error) {
	var payload []byte

	err := db.QueryRowContext(
		ctx,
		"SELECT payload FROM read_model_ops_bookings WHERE payload::jsonb -> 'tickets' ? $1",
		ticketID,
	).Scan(&payload)
	if err != nil {
		return entities.OpsBooking{}, err
	}

	return r.unmarshalReadModelFromDB(payload)
}

func (r OpsBookingReadModel) findReadModelByBookingID(
	ctx context.Context,
	bookingID string,
	db dbExecutor,
) (entities.OpsBooking, error) {
	var payload []byte

	err := db.QueryRowContext(
		ctx,
		"SELECT payload FROM read_model_ops_bookings WHERE booking_id = $1",
		bookingID,
	).Scan(&payload)
	if err != nil {
		return entities.OpsBooking{}, err
	}

	return r.unmarshalReadModelFromDB(payload)
}

func (r OpsBookingReadModel) unmarshalReadModelFromDB(payload []byte) (entities.OpsBooking, error) {
	var dbReadModel entities.OpsBooking
	if err := json.Unmarshal(payload, &dbReadModel); err != nil {
		return entities.OpsBooking{}, err
	}

	if dbReadModel.Tickets == nil {
		dbReadModel.Tickets = map[string]entities.OpsTicket{}
	}

	return dbReadModel, nil
}

type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func updateInTx(
	ctx context.Context,
	db *sqlx.DB,
	isolation sql.IsolationLevel,
	fn func(ctx context.Context, tx *sqlx.Tx) error,
) (err error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: isolation})
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
			return
		}

		err = tx.Commit()
	}()

	return fn(ctx, tx)
}
