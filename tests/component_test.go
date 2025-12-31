package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tickets/adapters"
	"tickets/entities"
	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/service"
)

func TestComponent(t *testing.T) {
	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	spreadsheetsAPI := &adapters.SpreadsheetsAPIStub{}
	receiptsService := &adapters.ReceiptsServiceStub{}

	go func() {
		svc := service.New(
			redisClient,
			spreadsheetsAPI,
			receiptsService,
		)
		assert.NoError(t, svc.Run(ctx))
	}()

	waitForHttpServer(t)

	ticket := ticketsHttp.TicketStatusRequest{
		TicketID: uuid.NewString(),
		Status:   "confirmed",
		Price: entities.Money{
			Amount:   "50.30",
			Currency: "GBP",
		},
		CustomerEmail: "email@example.com",
	}

	sendTicketsStatus(
		t,
		ticketsHttp.TicketsStatusRequest{
			Tickets: []ticketsHttp.TicketStatusRequest{ticket},
		},
	)

	assertReceiptForTicketIssued(t, receiptsService, ticket)
	assertRowToSheetAdded(
		t,
		spreadsheetsAPI,
		ticket,
		"tickets-to-print",
	)

	ticket.Status = "canceled"
	sendTicketsStatus(t, ticketsHttp.TicketsStatusRequest{
		Tickets: []ticketsHttp.TicketStatusRequest{ticket},
	})

	assertRowToSheetAdded(
		t,
		spreadsheetsAPI,
		ticket,
		"tickets-to-refund",
	)
}

func assertRowToSheetAdded(t *testing.T, spreadsheetsAPI *adapters.SpreadsheetsAPIStub, ticket ticketsHttp.TicketStatusRequest, sheetName string) bool {
	t.Helper()

	return assert.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			rows, ok := spreadsheetsAPI.Rows[sheetName]
			if !assert.True(t, ok, "sheet %s not found", sheetName) {
				return
			}

			var ticketRow []string

			for _, row := range rows {
				for _, col := range row {
					if col == ticket.TicketID {
						ticketRow = row
						break
					}
				}
			}
			if !assert.NotEmpty(t, ticketRow, "ticket row not found in sheet %s", sheetName) {
				return
			}

			expectedRow := []string{
				ticket.TicketID,
				ticket.CustomerEmail,
				ticket.Price.Amount,
				ticket.Price.Currency,
			}

			assert.Equal(
				t,
				expectedRow,
				ticketRow,
			)
		},
		10*time.Second,
		100*time.Millisecond,
	)
}

func assertReceiptForTicketIssued(t *testing.T, receiptsService *adapters.ReceiptsServiceStub, ticket ticketsHttp.TicketStatusRequest) {
	t.Helper()

	parentT := t

	assert.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			issuedReceipts := len(receiptsService.IssuedReceipts)
			parentT.Log("issued receipts", issuedReceipts)

			assert.Greater(t, issuedReceipts, 0, "no receipts issued")
		},
		10*time.Second,
		100*time.Millisecond,
	)

	var receipt entities.IssueReceiptRequest
	var ok bool
	for _, issuedReceipt := range receiptsService.IssuedReceipts {
		if issuedReceipt.TicketID != ticket.TicketID {
			continue
		}
		receipt = issuedReceipt
		ok = true
		break
	}
	require.Truef(t, ok, "receipt for ticket %s not found", ticket.TicketID)

	assert.Equal(t, ticket.TicketID, receipt.TicketID)
	assert.Equal(t, ticket.Price.Amount, receipt.Price.Amount)
	assert.Equal(t, ticket.Price.Currency, receipt.Price.Currency)
}

func sendTicketsStatus(t *testing.T, req ticketsHttp.TicketsStatusRequest) {
	t.Helper()

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	correlationID := shortuuid.New()

	httpReq, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/tickets-status",
		bytes.NewBuffer(payload),
	)
	require.NoError(t, err)

	httpReq.Header.Set("Correlation-ID", correlationID)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func waitForHttpServer(t *testing.T) {
	t.Helper()

	require.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			resp, err := http.Get("http://localhost:8080/health")
			if !assert.NoError(t, err) {
				return
			}
			defer resp.Body.Close()

			if assert.Less(t, resp.StatusCode, 300, "API not ready, http status: %d", resp.StatusCode) {
				return
			}
		},
		time.Second*10,
		time.Millisecond*50,
	)
}
