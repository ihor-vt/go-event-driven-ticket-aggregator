package event

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"

	"tickets/entities"
)

type Handler struct {
	deadNationAPI     DeadNationAPI
	spreadsheetsAPI   SpreadsheetsAPI
	receiptsService   ReceiptsService
	filesAPI          FilesAPI
	ticketsRepository TicketsRepository
	showsRepository   ShowsRepository
	eventBus          *cqrs.EventBus
}

func NewHandler(
	deadNationAPI DeadNationAPI,
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	filesAPI FilesAPI,
	ticketsRepository TicketsRepository,
	showsRepository ShowsRepository,
	eventBus *cqrs.EventBus,
) Handler {
	if eventBus == nil {
		panic("missing eventBus")
	}
	if deadNationAPI == nil {
		panic("missing deadNationAPI")
	}
	if spreadsheetsAPI == nil {
		panic("missing spreadsheetsAPI")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}
	if filesAPI == nil {
		panic("missing filesAPI")
	}
	if ticketsRepository == nil {
		panic("missing ticketsRepository")
	}
	if showsRepository == nil {
		panic("missing showsRepository")
	}
	if eventBus == nil {
		panic("missing eventBus")
	}

	return Handler{
		deadNationAPI:     deadNationAPI,
		spreadsheetsAPI:   spreadsheetsAPI,
		receiptsService:   receiptsService,
		filesAPI:          filesAPI,
		ticketsRepository: ticketsRepository,
		showsRepository:   showsRepository,
		eventBus:          eventBus,
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) error
}

type FilesAPI interface {
	UploadFile(ctx context.Context, fileID string, fileContent string) error
}

type TicketsRepository interface {
	Add(ctx context.Context, ticket entities.Ticket) error
	Remove(ctx context.Context, ticketID string) error
}

type ShowsRepository interface {
	ShowByID(ctx context.Context, showID uuid.UUID) (entities.Show, error)
}

type DeadNationAPI interface {
	BookInDeadNation(ctx context.Context, request entities.DeadNationBooking) error
}
